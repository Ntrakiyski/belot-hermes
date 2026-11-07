package belot

import (
	"errors"
	"fmt"
	"io"
	"sort"
)

// RulesetVersion is pinned to every game for later dispute/replay work.
// Bulgarian rules fixtures and remaining variant limits are in docs/06-rule-verification.md.
const RulesetVersion = "bg-0.3.0"

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
	Contract          Contract `json:"contract"`
	Bidder            int      `json:"bidder"`
	Multiplier        int      `json:"multiplier"`
	RawPoints         [2]int   `json:"rawPoints"`
	Tricks            [2]int   `json:"tricks"`
	Capot             int      `json:"capot"` // team index that won all tricks, or -1
	Succeeded         bool     `json:"succeeded"`
	Awarded           [2]int   `json:"awarded"`
	Notes             []string `json:"notes"`
	AllPassed         bool     `json:"allPassed"`
	CardPoints        [2]int   `json:"cardPoints"`
	CapotPoints       [2]int   `json:"capotPoints"`
	DeclarationPoints [2]int   `json:"declarationPoints"`
	BelotPoints       [2]int   `json:"belotPoints"`
	Hanging           bool     `json:"hanging"`
	Carried           int      `json:"carried"`
	CarryAwarded      int      `json:"carryAwarded"`
}

type Game struct {
	Dealer      int
	Phase       Phase
	Hands       [4][]Card
	rest        []Card // undealt 12 cards until bidding ends
	Turn        int
	Contract    Contract // "" until someone bids
	Bidder      int
	Mult        int // 1, 2 (contra), 4 (recontra)
	passes      int
	Trick       []Played
	LastTrick   []Played
	Captured    [2]int
	Tricks      [2]int
	Scores      [2]int // match totals
	Target      int
	Result      *DealResult
	Winner      int // winning team when Phase==PhaseMatch
	Log         []string
	Seq         int
	Pending     int
	BelotPoints [2]int
	declared    []declaredCombination
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
	g.Captured, g.Tricks, g.BelotPoints = [2]int{}, [2]int{}, [2]int{}
	g.declared = nil
	for i := 0; i < 4; i++ {
		seat := (g.Dealer + 1 + i) % 4
		g.Hands[seat] = append([]Card(nil), deck[i*3:i*3+3]...)
	}
	for i := 0; i < 4; i++ {
		seat := (g.Dealer + 1 + i) % 4
		g.Hands[seat] = append(g.Hands[seat], deck[12+i*2:14+i*2]...)
	}
	g.rest = append([]Card(nil), deck[20:]...)
	g.Turn = (g.Dealer + 1) % 4
	g.logf("New deal. Seat %d deals; seat %d bids first.", g.Dealer, g.Turn)
	g.sortHands()
	return nil
}

var ErrNotYourTurn = errors.New("Изчакайте своя ред.")

// Bid handles action = "pass" | "bid" | "double" | "redouble".
func (g *Game) Bid(seat int, action string, c Contract) error {
	if g.Phase != PhaseBidding {
		return errors.New("Наддаването вече приключи.")
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
			g.Result = &DealResult{AllPassed: true, Capot: -1, Carried: g.Pending, Notes: []string{"Всички обявиха пас: събираме картите и раздава следващият играч."}}
			g.Phase = PhaseDealOver
			return nil
		}
		if g.Contract != "" && g.passes == 3 {
			g.finishBidding()
			return nil
		}
	case "bid":
		if !c.Valid() {
			return errors.New("Непозната обява.")
		}
		if g.Contract != "" && contractOrder[c] <= contractOrder[g.Contract] {
			return errors.New("Обявата трябва да е по-висока от текущата.")
		}
		g.Contract, g.Bidder, g.Mult, g.passes = c, seat, 1, 0
		g.logf("Seat %d bids %s.", seat, c)
	case "double":
		if g.Contract == "" || g.Mult != 1 || Team(seat) == Team(g.Bidder) {
			return errors.New("В момента не може да обявите контра.")
		}
		g.Mult, g.passes = 2, 0
		g.logf("Seat %d doubles (contra).", seat)
	case "redouble":
		if g.Contract == "" || g.Mult != 2 || Team(seat) != Team(g.Bidder) {
			return errors.New("В момента не може да обявите реконтра.")
		}
		g.Mult, g.passes = 4, 0
		g.logf("Seat %d redoubles (recontra).", seat)
	default:
		return errors.New("Непознато действие при наддаване.")
	}
	g.Turn = (g.Turn + 1) % 4
	return nil
}

func (g *Game) finishBidding() {
	for i := 0; i < 4; i++ {
		seat := (g.Dealer + 1 + i) % 4
		g.Hands[seat] = append(g.Hands[seat], g.rest[i*3:i*3+3]...)
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
		return append([]Card(nil), hand...), "Вие сте на ход: изберете произволна карта."
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
			return follow, "Отговорете на поисканата боя."
		}
		return append([]Card(nil), hand...), "Нямате поисканата боя: може да играете всяка карта."
	case AllTrump:
		if len(follow) > 0 {
			if h := higher(follow); len(h) > 0 {
				return h, "Отговорете на боята и качете, ако можете."
			}
			return follow, "Отговорете на боята: нямате по-висока карта."
		}
		return append([]Card(nil), hand...), "Нямате поисканата боя: може да играете всяка карта."
	}
	trump, _ := g.Contract.trumpSuit()
	if led == trump {
		if len(follow) > 0 {
			if h := higher(follow); len(h) > 0 {
				return h, "Поискан е коз: качете с по-висок коз, ако можете."
			}
			return follow, "Поискан е коз: играйте коз."
		}
		return append([]Card(nil), hand...), "Нямате коз: може да играете всяка карта."
	}
	if len(follow) > 0 {
		return follow, "Отговорете на поисканата боя."
	}
	if Team(best.Seat) == Team(seat) {
		return append([]Card(nil), hand...), "Партньорът ви печели: може да играете всяка карта."
	}
	var trumps []Card
	for _, c := range hand {
		if c.Suit == trump {
			trumps = append(trumps, c)
		}
	}
	if len(trumps) == 0 {
		return append([]Card(nil), hand...), "Нямате поисканата боя и коз: може да играете всяка карта."
	}
	if g.isTrump(best.Card) {
		if h := higher(trumps); len(h) > 0 {
			return h, "Трябва да надцакате."
		}
		return append([]Card(nil), hand...), "Не можете да надцакате: може да играете всяка карта."
	}
	return trumps, "Нямате поисканата боя: трябва да цакате."
}

func (g *Game) Play(seat int, card Card) error {
	return g.PlayDeclared(seat, card, nil, false)
}

// PlayDeclared validates the card and optional announcements before changing any state.
func (g *Game) PlayDeclared(seat int, card Card, declarationIDs []string, announceBelot bool) error {
	if g.Phase != PhasePlaying {
		return errors.New("В момента не се разиграва.")
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
			return errors.New("Тази карта не е в ръката ви.")
		}
		return fmt.Errorf("Недопустима карта: %s", reason)
	}
	selected, err := g.selectedDeclarations(seat, declarationIDs)
	if err != nil {
		return err
	}
	if announceBelot && !g.CanAnnounceBelot(seat, card) {
		return errors.New("Белот може да се обяви само при допустимото изиграване на първия поп или дама.")
	}
	for _, d := range selected {
		g.declared = append(g.declared, declaredCombination{Seat: seat, Declaration: d})
		g.logf("Seat %d announces %s (+%d).", seat, d.Name, d.Points)
	}
	if announceBelot {
		g.BelotPoints[Team(seat)] += 20
		g.declared = append(g.declared, declaredCombination{Seat: seat, Declaration: Declaration{Kind: "belot", Name: "Белот", Points: 20}})
		g.logf("Seat %d announces белот (+20).", seat)
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

func roundTotal(raw int) int { return (raw + 5) / 10 }

// Bulgarian split scoring uses paired thresholds; rounding teams independently loses
// a point at the shared 6 (suit) or shared 4 (all-trump) boundary.
func roundedTeams(raw [2]int, contract Contract) [2]int {
	offset := 3
	if contract == AllTrump || contract == NoTrump {
		offset = 5
	}
	points := [2]int{(raw[0] + offset) / 10, (raw[1] + offset) / 10}
	boundary := 6
	if contract == AllTrump {
		boundary = 4
	}
	if contract != NoTrump && raw[0]%10 == boundary && raw[1]%10 == boundary && raw[0] != raw[1] {
		lower := 0
		if raw[1] < raw[0] {
			lower = 1
		}
		points[lower]++
	}
	return points
}

func (g *Game) finishDeal() {
	r := &DealResult{Contract: g.Contract, Bidder: g.Bidder, Multiplier: g.Mult, Capot: -1, Tricks: g.Tricks, DeclarationPoints: g.DeclarationPoints(), BelotPoints: g.BelotPoints}
	raw := g.Captured
	if g.Contract == NoTrump {
		raw[0] *= 2
		raw[1] *= 2
		r.Notes = append(r.Notes, "Без коз: точките от карти и последно 10 се удвояват.")
	}
	r.CardPoints = raw
	for team := 0; team < 2; team++ {
		raw[team] += r.DeclarationPoints[team] + r.BelotPoints[team]
		if g.Tricks[team] == 8 {
			r.Capot = team
			r.CapotPoints[team] = 90
			raw[team] += 90
			r.Notes = append(r.Notes, fmt.Sprintf("Отбор %d взема всички 8 взятки: +90 за капо.", team+1))
		}
	}
	r.RawPoints = raw
	bt := Team(g.Bidder)
	dt := 1 - bt
	r.Succeeded = raw[bt] > raw[dt]
	r.Hanging = raw[bt] == raw[dt]
	pot := roundTotal(raw[0] + raw[1])
	split := roundedTeams(raw, g.Contract)
	switch {
	case r.Hanging && g.Mult > 1:
		g.Pending += pot * g.Mult
		r.Notes = append(r.Notes, fmt.Sprintf("Висяща игра при контра/реконтра: %d точки се пренасят.", pot*g.Mult))
	case r.Hanging:
		r.Awarded[dt] = split[dt]
		g.Pending += pot - split[dt]
		r.Notes = append(r.Notes, fmt.Sprintf("Висяща игра: защитата записва %d, а %d точки се пренасят.", split[dt], pot-split[dt]))
	case g.Mult > 1:
		winner := dt
		if r.Succeeded {
			winner = bt
		}
		r.Awarded[winner] = pot * g.Mult
		r.Notes = append(r.Notes, fmt.Sprintf("Контра/реконтра: отборът с повече точки записва целия сбор ×%d.", g.Mult))
	case r.Succeeded:
		r.Awarded = split
		r.Notes = append(r.Notes, "Изкарана игра: всеки отбор записва спечелените точки след закръгляне.")
	default:
		r.Awarded[dt] = pot
		r.Notes = append(r.Notes, "Вкарана игра: защитата записва целия сбор.")
	}
	if !r.Hanging && g.Pending > 0 {
		winner := dt
		if r.Succeeded {
			winner = bt
		}
		r.CarryAwarded = g.Pending
		r.Awarded[winner] += g.Pending
		g.Pending = 0
		r.Notes = append(r.Notes, fmt.Sprintf("Победителят получава и %d пренесени точки.", r.CarryAwarded))
	}
	r.Carried = g.Pending
	g.Scores[0] += r.Awarded[0]
	g.Scores[1] += r.Awarded[1]
	g.Result = r
	g.Phase = PhaseDealOver
	g.logf("Deal over: team0 +%d, team1 +%d.", r.Awarded[0], r.Awarded[1])
	if r.Capot < 0 {
		g.checkMatch()
	} else {
		r.Notes = append(r.Notes, "С капо не се излиза: следва още едно раздаване.")
	}
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
		return errors.New("Раздаването още не е приключило.")
	}
	g.Dealer = (g.Dealer + 1) % 4
	return g.StartDeal(src)
}

// WinningSeat returns the winner of the public trick so far, or -1 before its lead.
func (g *Game) WinningSeat() int {
	if len(g.Trick) == 0 {
		return -1
	}
	return g.currentWinner().Seat
}

// Beats reports whether candidate beats current for the led suit under this contract.
func (g *Game) Beats(led Suit, current, candidate Card) bool {
	return g.trickBeats(led, current, candidate)
}
