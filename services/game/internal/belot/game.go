package belot

import (
	"errors"
	"fmt"
	"io"
	"sort"
)

// RulesetVersion is pinned to every game for later dispute/replay work.
// "mvp-0" = playable baseline; see docs/06-rule-verification.md for what is NOT yet verified
// (announcements/tierce/belot, hanging hands, exact rounding).
const RulesetVersion = "mvp-0"

type Contract string

const (
	Clubs    Contract = "clubs"
	Diamonds Contract = "diamonds"
	Hearts   Contract = "hearts"
	Spades   Contract = "spades"
	NoTrump  Contract = "notrump"
	AllTrump Contract = "alltrump"
)

// contractOrder is the bid ladder: clubs < diamonds < hearts < spades < no trump < all trump.
var contractOrder = map[Contract]int{Clubs: 1, Diamonds: 2, Hearts: 3, Spades: 4, NoTrump: 5, AllTrump: 6}

func (c Contract) Valid() bool { _, ok := contractOrder[c]; return ok }

func (c Contract) trumpSuit() (Suit, bool) {
	switch c {
	case Clubs, Diamonds, Hearts, Spades:
		return Suit(c), true
	}
	return "", false
}

type Phase string

const (
	PhaseBidding  Phase = "bidding"
	PhasePlaying  Phase = "playing"
	PhaseDealOver Phase = "dealover"
	PhaseMatch    Phase = "matchover"
)

type Played struct {
	Seat int  `json:"seat"`
	Card Card `json:"card"`
}

// DealResult is the explained score of a finished deal.
type DealResult struct {
	Contract   Contract `json:"contract"`
	Bidder     int      `json:"bidder"`
	Multiplier int      `json:"multiplier"`
	RawPoints  [2]int   `json:"rawPoints"`
	Tricks     [2]int   `json:"tricks"`
	Capot      int      `json:"capot"` // team index that won all tricks, or -1
	Succeeded  bool     `json:"succeeded"`
	Awarded    [2]int   `json:"awarded"`
	Notes      []string `json:"notes"`
	AllPassed  bool     `json:"allPassed"`
}

type Game struct {
	Dealer    int
	Phase     Phase
	Hands     [4][]Card
	rest      []Card // undealt 12 cards until bidding ends
	Turn      int
	Contract  Contract // "" until someone bids
	Bidder    int
	Mult      int // 1, 2 (contra), 4 (recontra)
	passes    int
	Trick     []Played
	LastTrick []Played
	Captured  [2]int
	Tricks    [2]int
	Scores    [2]int // match totals
	Target    int
	Result    *DealResult
	Winner    int // winning team when Phase==PhaseMatch
	Log       []string
	Seq       int
}

func Team(seat int) int { return seat % 2 }

func NewGame(dealer int, target int) *Game {
	if target == 0 {
		target = 151
	}
	return &Game{Dealer: dealer, Target: target, Winner: -1}
}

func (g *Game) logf(format string, a ...any) {
	g.Seq++
	g.Log = append(g.Log, fmt.Sprintf(format, a...))
	if len(g.Log) > 60 {
		g.Log = g.Log[len(g.Log)-60:]
	}
}

// StartDeal shuffles (or uses the supplied RNG) and deals 5 cards each; 12 remain for after bidding.
func (g *Game) StartDeal(src io.Reader) error {
	deck, err := Shuffle(src)
	if err != nil {
		return err
	}
	return g.StartDealWith(deck)
}

func (g *Game) StartDealWith(deck []Card) error {
	if err := ValidateDeck(deck); err != nil {
		return err
	}
	g.Phase = PhaseBidding
	g.Contract, g.Mult, g.passes, g.Bidder = "", 1, 0, -1
	g.Trick, g.LastTrick, g.Result = nil, nil, nil
	g.Captured, g.Tricks = [2]int{}, [2]int{}
	for s := 0; s < 4; s++ {
		g.Hands[s] = append([]Card(nil), deck[s*5:s*5+5]...)
	}
	g.rest = append([]Card(nil), deck[20:]...)
	g.Turn = (g.Dealer + 1) % 4
	g.logf("New deal. Seat %d deals; seat %d bids first.", g.Dealer, g.Turn)
	g.sortHands()
	return nil
}

var ErrNotYourTurn = errors.New("not your turn")

// Bid handles action = "pass" | "bid" | "double" | "redouble".
func (g *Game) Bid(seat int, action string, c Contract) error {
	if g.Phase != PhaseBidding {
		return errors.New("not in bidding phase")
	}
	if seat != g.Turn {
		return ErrNotYourTurn
	}
	switch action {
	case "pass":
		g.passes++
		g.logf("Seat %d passes.", seat)
		if g.Contract == "" && g.passes == 4 {
			g.logf("Everyone passed: redeal.")
			g.Result = &DealResult{AllPassed: true, Capot: -1, Notes: []string{"All four players passed; the deal is thrown in and the next dealer deals."}}
			g.Phase = PhaseDealOver
			return nil
		}
		if g.Contract != "" && g.passes == 3 {
			g.finishBidding()
			return nil
		}
	case "bid":
		if !c.Valid() {
			return errors.New("unknown contract")
		}
		if g.Contract != "" && contractOrder[c] <= contractOrder[g.Contract] {
			return errors.New("bid must be higher than the current contract")
		}
		g.Contract, g.Bidder, g.Mult, g.passes = c, seat, 1, 0
		g.logf("Seat %d bids %s.", seat, c)
	case "double":
		if g.Contract == "" || g.Mult != 1 || Team(seat) == Team(g.Bidder) {
			return errors.New("cannot double now")
		}
		g.Mult, g.passes = 2, 0
		g.logf("Seat %d doubles (contra).", seat)
	case "redouble":
		if g.Contract == "" || g.Mult != 2 || Team(seat) != Team(g.Bidder) {
			return errors.New("cannot redouble now")
		}
		g.Mult = 4
		g.logf("Seat %d redoubles (recontra).", seat)
		g.finishBidding()
		return nil
	default:
		return errors.New("unknown bidding action")
	}
	g.Turn = (g.Turn + 1) % 4
	return nil
}

func (g *Game) finishBidding() {
	for s := 0; s < 4; s++ {
		g.Hands[s] = append(g.Hands[s], g.rest[s*3:s*3+3]...)
	}
	g.rest = nil
	g.Phase = PhasePlaying
	g.Turn = (g.Dealer + 1) % 4
	g.sortHands()
	g.logf("Contract: %s by seat %d (x%d). Seat %d leads.", g.Contract, g.Bidder, g.Mult, g.Turn)
}

func (g *Game) sortHands() {
	for s := 0; s < 4; s++ {
		h := g.Hands[s]
		sort.Slice(h, func(i, j int) bool {
			si, sj := suitIndex(h[i].Suit), suitIndex(h[j].Suit)
			if si != sj {
				return si < sj
			}
			return rankIndex(h[i].Rank, nil) < rankIndex(h[j].Rank, nil)
		})
	}
}

func suitIndex(s Suit) int {
	for i, x := range Suits {
		if x == s {
			return i
		}
	}
	return -1
}

var normalOrder = []Rank{"7", "8", "9", "J", "Q", "K", "10", "A"}
var trumpOrder = []Rank{"7", "8", "Q", "K", "10", "A", "9", "J"}

func rankIndex(r Rank, order []Rank) int {
	if order == nil {
		order = normalOrder
	}
	for i, x := range order {
		if x == r {
			return i
		}
	}
	return -1
}

// strength of a card inside its own suit under the contract.
func (g *Game) strength(c Card) int {
	if g.Contract == AllTrump {
		return rankIndex(c.Rank, trumpOrder)
	}
	if t, ok := g.Contract.trumpSuit(); ok && c.Suit == t {
		return rankIndex(c.Rank, trumpOrder)
	}
	return rankIndex(c.Rank, normalOrder)
}

func (g *Game) isTrump(c Card) bool {
	if t, ok := g.Contract.trumpSuit(); ok {
		return c.Suit == t
	}
	return false
}

// CardPoints returns the card's point value under the current contract (before no-trump doubling).
func (g *Game) CardPoints(c Card) int {
	trumpLike := g.Contract == AllTrump || g.isTrump(c)
	if trumpLike {
		switch c.Rank {
		case "J":
			return 20
		case "9":
			return 14
		case "A":
			return 11
		case "10":
			return 10
		case "K":
			return 4
		case "Q":
			return 3
		}
		return 0
	}
	switch c.Rank {
	case "A":
		return 11
	case "10":
		return 10
	case "K":
		return 4
	case "Q":
		return 3
	case "J":
		return 2
	}
	return 0
}

// trickBeats reports whether b beats a for the given led suit.
func (g *Game) trickBeats(led Suit, a, b Card) bool {
	at, bt := g.isTrump(a), g.isTrump(b)
	if at != bt {
		return bt
	}
	if a.Suit != b.Suit {
		return false // b did not follow and is not a trump
	}
	return g.strength(b) > g.strength(a)
}

func (g *Game) currentWinner() Played {
	led := g.Trick[0].Card.Suit
	best := g.Trick[0]
	for _, p := range g.Trick[1:] {
		if g.trickBeats(led, best.Card, p.Card) {
			best = p
		}
	}
	return best
}

// LegalCards returns the legal cards for seat and a short reason string.
func (g *Game) LegalCards(seat int) ([]Card, string) {
	if g.Phase != PhasePlaying || seat != g.Turn {
		return nil, ""
	}
	hand := g.Hands[seat]
	if len(g.Trick) == 0 {
		return append([]Card(nil), hand...), "You lead: any card."
	}
	led := g.Trick[0].Card.Suit
	var follow []Card
	for _, c := range hand {
		if c.Suit == led {
			follow = append(follow, c)
		}
	}
	best := g.currentWinner()
	higher := func(cards []Card) []Card {
		var out []Card
		for _, c := range cards {
			if g.trickBeats(led, best.Card, c) {
				out = append(out, c)
			}
		}
		return out
	}
	switch g.Contract {
	case NoTrump:
		if len(follow) > 0 {
			return follow, "Follow the led suit."
		}
		return append([]Card(nil), hand...), "No card of the led suit: play anything."
	case AllTrump:
		if len(follow) > 0 {
			if h := higher(follow); len(h) > 0 {
				return h, "Follow suit and beat the highest card if you can."
			}
			return follow, "Follow suit (you cannot beat the highest card)."
		}
		return append([]Card(nil), hand...), "No card of the led suit: play anything."
	}
	trump, _ := g.Contract.trumpSuit()
	if led == trump {
		if len(follow) > 0 {
			if h := higher(follow); len(h) > 0 {
				return h, "Trump led: you must beat the highest trump if you can."
			}
			return follow, "Trump led: follow with a trump."
		}
		return append([]Card(nil), hand...), "No trumps: play anything."
	}
	if len(follow) > 0 {
		return follow, "Follow the led suit."
	}
	if Team(best.Seat) == Team(seat) {
		return append([]Card(nil), hand...), "Your partner is winning: play anything."
	}
	var trumps []Card
	for _, c := range hand {
		if c.Suit == trump {
			trumps = append(trumps, c)
		}
	}
	if len(trumps) == 0 {
		return append([]Card(nil), hand...), "No led suit and no trumps: play anything."
	}
	if g.isTrump(best.Card) {
		if h := higher(trumps); len(h) > 0 {
			return h, "You must over-trump."
		}
		return append([]Card(nil), hand...), "You cannot over-trump: play anything."
	}
	return trumps, "You cannot follow suit: you must trump."
}

func (g *Game) Play(seat int, card Card) error {
	if g.Phase != PhasePlaying {
		return errors.New("not in playing phase")
	}
	if seat != g.Turn {
		return ErrNotYourTurn
	}
	legal, reason := g.LegalCards(seat)
	ok := false
	for _, c := range legal {
		if c == card {
			ok = true
		}
	}
	if !ok {
		in := false
		for _, c := range g.Hands[seat] {
			if c == card {
				in = true
			}
		}
		if !in {
			return errors.New("that card is not in your hand")
		}
		return fmt.Errorf("illegal card: %s", reason)
	}
	h := g.Hands[seat]
	for i, c := range h {
		if c == card {
			g.Hands[seat] = append(h[:i:i], h[i+1:]...)
			break
		}
	}
	g.Trick = append(g.Trick, Played{Seat: seat, Card: card})
	g.logf("Seat %d plays %s %s.", seat, card.Rank, card.Suit)
	if len(g.Trick) < 4 {
		g.Turn = (g.Turn + 1) % 4
		return nil
	}
	w := g.currentWinner()
	pts := 0
	for _, p := range g.Trick {
		pts += g.CardPoints(p.Card)
	}
	t := Team(w.Seat)
	g.Tricks[t]++
	last := len(g.Hands[w.Seat]) == 0
	if last {
		pts += 10
	}
	g.Captured[t] += pts
	g.logf("Seat %d wins the trick (+%d).", w.Seat, pts)
	g.LastTrick, g.Trick = g.Trick, nil
	g.Turn = w.Seat
	if last {
		g.finishDeal()
	}
	return nil
}

func roundTens(raw int) int { return (raw + 4) / 10 }

func (g *Game) finishDeal() {
	r := &DealResult{Contract: g.Contract, Bidder: g.Bidder, Multiplier: g.Mult, Capot: -1, Tricks: g.Tricks}
	raw := g.Captured
	mult := 1
	if g.Contract == NoTrump {
		mult = 2
		r.Notes = append(r.Notes, "No trumps: card points are doubled.")
	}
	raw[0], raw[1] = raw[0]*mult, raw[1]*mult
	total := raw[0] + raw[1]
	for t := 0; t < 2; t++ {
		if g.Tricks[t] == 8 {
			r.Capot = t
			raw[t] += 90
			r.Notes = append(r.Notes, fmt.Sprintf("Team %d won all 8 tricks (capot): +90.", t))
		}
	}
	r.RawPoints = raw
	bt := Team(g.Bidder)
	dt := 1 - bt
	succeeded := raw[bt]-boolInt(r.Capot == bt)*90 > total/2
	r.Succeeded = succeeded
	pot := roundTens(raw[0]) + roundTens(raw[1])
	if r.Capot >= 0 {
		pot = roundTens(raw[0]) + roundTens(raw[1])
	}
	switch {
	case g.Mult > 1 && succeeded:
		r.Awarded[bt] = pot * g.Mult
		r.Notes = append(r.Notes, fmt.Sprintf("Contract made against a double: bidding team takes the whole pot x%d.", g.Mult))
	case g.Mult > 1:
		r.Awarded[dt] = pot * g.Mult
		r.Notes = append(r.Notes, fmt.Sprintf("Contract failed against a double: defenders take the whole pot x%d.", g.Mult))
	case succeeded:
		r.Awarded[bt], r.Awarded[dt] = roundTens(raw[bt]), roundTens(raw[dt])
		r.Notes = append(r.Notes, "Contract made: each team scores its own points.")
	default:
		r.Awarded[dt] = pot
		r.Notes = append(r.Notes, "Contract failed (needs MORE than half): defenders take all the points.")
	}
	g.Scores[0] += r.Awarded[0]
	g.Scores[1] += r.Awarded[1]
	g.Result = r
	g.Phase = PhaseDealOver
	g.logf("Deal over: team0 +%d, team1 +%d.", r.Awarded[0], r.Awarded[1])
	g.checkMatch()
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (g *Game) checkMatch() {
	a, b := g.Scores[0], g.Scores[1]
	if (a >= g.Target || b >= g.Target) && a != b {
		g.Phase = PhaseMatch
		if a > b {
			g.Winner = 0
		} else {
			g.Winner = 1
		}
		g.logf("Match over: team %d wins %d-%d.", g.Winner, a, b)
	}
}

// NextDeal rotates the dealer and deals again (after a deal result has been acknowledged).
func (g *Game) NextDeal(src io.Reader) error {
	if g.Phase != PhaseDealOver {
		return errors.New("deal is not finished")
	}
	g.Dealer = (g.Dealer + 1) % 4
	return g.StartDeal(src)
}
