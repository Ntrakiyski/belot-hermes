package main

import (
	"sort"
	"time"

	"github.com/Ntrakiyski/belot-hermes/services/game/internal/belot"
)

const (
	botBidDelay  = 800 * time.Millisecond
	botPlayDelay = 650 * time.Millisecond
	trickHold    = 1800 * time.Millisecond
)

type bidEvent struct {
	Seat     int            `json:"seat"`
	Action   string         `json:"action"`
	Contract belot.Contract `json:"contract"`
}

type dealSummary struct {
	Number    int            `json:"number"`
	Contract  belot.Contract `json:"contract"`
	Awarded   [2]int         `json:"awarded"`
	Scores    [2]int         `json:"scores"`
	AllPassed bool           `json:"allPassed"`
	Hanging   bool           `json:"hanging"`
}

func (s *store) run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for now := range ticker.C {
		s.tick(now)
	}
}

func (s *store) tick(now time.Time) {
	s.mu.Lock()
	tables := make([]*table, 0, len(s.tables))
	for _, t := range s.tables {
		tables = append(tables, t)
	}
	s.mu.Unlock()
	for _, t := range tables {
		t.mu.Lock()
		t.advance(now)
		t.mu.Unlock()
	}
}

func (t *table) schedule(now time.Time) {
	t.NextAction = time.Time{}
	if t.Game == nil || !t.CollectUntil.IsZero() {
		return
	}
	g := t.Game
	if g.Phase != belot.PhaseBidding && g.Phase != belot.PhasePlaying {
		return
	}
	if s := t.Seats[g.Turn]; s != nil && s.Bot {
		delay := botPlayDelay
		if g.Phase == belot.PhaseBidding {
			delay = botBidDelay
		}
		t.NextAction = now.Add(delay)
	}
}

func (t *table) progress(now time.Time, beforeTricks, beforePoints int) {
	t.Seq++
	g := t.Game
	if g.Tricks[0]+g.Tricks[1] > beforeTricks {
		t.TrickWinner = g.Turn
		t.TrickPoints = g.Captured[0] + g.Captured[1] - beforePoints
		t.CollectUntil = now.Add(trickHold)
	}
	if r := g.Result; r != nil && (len(t.History) == 0 || t.History[len(t.History)-1].Number != t.DealNumber) {
		t.History = append(t.History, dealSummary{Number: t.DealNumber, Contract: g.Contract, Awarded: r.Awarded, Scores: g.Scores, AllPassed: r.AllPassed, Hanging: r.Hanging})
	}
	t.schedule(now)
}

// Each deadline advances one visible action, including when the ticker was delayed.
func (t *table) advance(now time.Time) {
	if !t.CollectUntil.IsZero() {
		if now.Before(t.CollectUntil) {
			return
		}
		t.CollectUntil = time.Time{}
		t.Seq++
		t.schedule(now)
		return
	}
	if t.NextAction.IsZero() || now.Before(t.NextAction) {
		return
	}
	g := t.Game
	t.NextAction = time.Time{}
	if g == nil || t.Seats[g.Turn] == nil || !t.Seats[g.Turn].Bot {
		return
	}
	beforeTricks := g.Tricks[0] + g.Tricks[1]
	beforePoints := g.Captured[0] + g.Captured[1]
	var err error
	if g.Phase == belot.PhaseBidding {
		action, contract := botBid(g)
		seat := g.Turn
		err = g.Bid(seat, action, contract)
		if err == nil {
			t.BidHistory = append(t.BidHistory, bidEvent{Seat: seat, Action: action, Contract: contract})
		}
	} else if g.Phase == belot.PhasePlaying {
		card, ok := botCard(g)
		if !ok {
			return
		}
		options := g.DeclarationOptions(g.Turn)
		// Announce a compatible set, choosing the higher-value combination on overlap.
		sort.SliceStable(options, func(i, j int) bool { return options[i].Points > options[j].Points })
		used := map[belot.Card]bool{}
		ids := []string{}
		for _, option := range options {
			overlap := false
			for _, c := range option.Cards {
				if used[c] {
					overlap = true
				}
			}
			if !overlap {
				ids = append(ids, option.ID)
				for _, c := range option.Cards {
					used[c] = true
				}
			}
		}
		err = g.PlayDeclared(g.Turn, card, ids, g.CanAnnounceBelot(g.Turn, card))
	} else {
		return
	}
	if err != nil {
		return
	}
	t.progress(now, beforeTricks, beforePoints)
}

var bidLadder = []belot.Contract{belot.Clubs, belot.Diamonds, belot.Hearts, belot.Spades, belot.NoTrump, belot.AllTrump}

// Practice decisions use this seat's cards and the public contract, never another hand.
func botBid(g *belot.Game) (string, belot.Contract) {
	current := -1
	for i, c := range bidLadder {
		if c == g.Contract {
			current = i
		}
	}
	bestMargin, best := 0, belot.Contract("")
	for i, c := range bidLadder {
		if i <= current {
			continue
		}
		probe := belot.Game{Contract: c}
		points, count := 0, 0
		for _, card := range g.Hands[g.Turn] {
			if c == belot.AllTrump || c == belot.NoTrump || card.Suit == belot.Suit(c) {
				points += probe.CardPoints(card)
				if card.Suit == belot.Suit(c) {
					count++
				}
			} else if card.Rank == "A" {
				points += 7
			}
		}
		threshold := 35
		if c == belot.NoTrump {
			threshold = 30
		} else if c == belot.AllTrump {
			threshold = 48
		} else {
			if count < 2 {
				continue
			}
			points += count * 4
		}
		if margin := points - threshold; margin > bestMargin {
			bestMargin, best = margin, c
		}
	}
	if best != "" {
		return "bid", best
	}
	return "pass", ""
}

func botCard(g *belot.Game) (belot.Card, bool) {
	legal, _ := g.LegalCards(g.Turn)
	if len(legal) == 0 {
		return belot.Card{}, false
	}
	sort.SliceStable(legal, func(i, j int) bool { return g.CardPoints(legal[i]) < g.CardPoints(legal[j]) })
	if len(g.Trick) == 0 {
		for _, c := range legal {
			if c.Rank == "A" && g.Contract != belot.AllTrump && c.Suit != belot.Suit(g.Contract) {
				return c, true
			}
		}
		return legal[0], true
	}
	winner := g.WinningSeat()
	if belot.Team(winner) == belot.Team(g.Turn) {
		return legal[0], true
	}
	var best belot.Card
	for _, p := range g.Trick {
		if p.Seat == winner {
			best = p.Card
		}
	}
	for _, c := range legal {
		if g.Beats(g.Trick[0].Card.Suit, best, c) {
			return c, true
		}
	}
	return legal[0], true
}
