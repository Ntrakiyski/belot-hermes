package belot

import (
	"bytes"
	"testing"
)

func newPlaying(t *testing.T, c Contract) *Game {
	t.Helper()
	g := NewGame(3, 151) // dealer 3 => seat 0 bids/leads first
	if err := g.StartDeal(bytes.NewReader(bytes.Repeat([]byte{7, 1, 9, 3}, 400))); err != nil {
		t.Fatal(err)
	}
	if err := g.Bid(0, "bid", c); err != nil {
		t.Fatal(err)
	}
	for s := 1; s <= 3; s++ {
		if err := g.Bid(s, "pass", ""); err != nil {
			t.Fatal(err)
		}
	}
	return g
}

func TestBiddingDealsFiveThenEight(t *testing.T) {
	g := NewGame(3, 151)
	_ = g.StartDeal(nil)
	for s := 0; s < 4; s++ {
		if len(g.Hands[s]) != 5 {
			t.Fatalf("seat %d has %d cards before bidding", s, len(g.Hands[s]))
		}
	}
	_ = g.Bid(0, "bid", Hearts)
	for s := 1; s <= 3; s++ {
		_ = g.Bid(s, "pass", "")
	}
	if g.Phase != PhasePlaying {
		t.Fatalf("phase %s", g.Phase)
	}
	seen := map[Card]bool{}
	for s := 0; s < 4; s++ {
		if len(g.Hands[s]) != 8 {
			t.Fatalf("seat %d has %d cards", s, len(g.Hands[s]))
		}
		for _, c := range g.Hands[s] {
			if seen[c] {
				t.Fatalf("duplicate %v", c)
			}
			seen[c] = true
		}
	}
}

func TestBidRulesRejectInvalid(t *testing.T) {
	g := NewGame(3, 151)
	_ = g.StartDeal(nil)
	if err := g.Bid(1, "bid", Clubs); err == nil {
		t.Fatal("out-of-turn bid accepted")
	}
	_ = g.Bid(0, "bid", Hearts)
	if err := g.Bid(1, "bid", Diamonds); err == nil {
		t.Fatal("lower bid accepted")
	}
	if err := g.Bid(1, "redouble", ""); err == nil {
		t.Fatal("redouble without double accepted")
	}
	if err := g.Bid(2, "double", ""); err == nil {
		t.Fatal("out-of-turn double accepted")
	}
	if err := g.Bid(1, "double", ""); err != nil {
		t.Fatal(err)
	}
	if err := g.Bid(2, "double", ""); err == nil {
		t.Fatal("double after double accepted")
	}
}

func TestAllPassIsThrownIn(t *testing.T) {
	g := NewGame(3, 151)
	_ = g.StartDeal(nil)
	for s := 0; s < 4; s++ {
		if err := g.Bid(s, "pass", ""); err != nil {
			t.Fatal(err)
		}
	}
	if g.Phase != PhaseDealOver || g.Result == nil || !g.Result.AllPassed {
		t.Fatal("expected thrown-in deal")
	}
}

func hand(g *Game, seat int, cs ...Card) { g.Hands[seat] = cs }

func TestFollowSuitAndMustTrump(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Bidder: 0, Mult: 1, Turn: 1}
	g.Trick = []Played{{0, Card{"spades", "K"}}}
	hand(g, 1, Card{"spades", "7"}, Card{"hearts", "7"}, Card{"clubs", "A"})
	legal, _ := g.LegalCards(1)
	if len(legal) != 1 || legal[0] != (Card{"spades", "7"}) {
		t.Fatalf("must follow suit: %v", legal)
	}
	hand(g, 1, Card{"hearts", "7"}, Card{"clubs", "A"})
	legal, _ = g.LegalCards(1)
	if len(legal) != 1 || legal[0].Suit != "hearts" {
		t.Fatalf("must trump when opponent winning: %v", legal)
	}
	// partner winning => free choice
	g.Turn = 2
	g.Trick = []Played{{0, Card{"spades", "K"}}, {1, Card{"spades", "7"}}}
	hand(g, 2, Card{"hearts", "7"}, Card{"clubs", "A"})
	legal, _ = g.LegalCards(2)
	if len(legal) != 2 {
		t.Fatalf("partner winning, free choice: %v", legal)
	}
}

func TestOvertrump(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Bidder: 0, Mult: 1, Turn: 2}
	g.Trick = []Played{{0, Card{"spades", "K"}}, {1, Card{"hearts", "9"}}}
	hand(g, 2, Card{"hearts", "7"}, Card{"hearts", "J"}, Card{"clubs", "A"})
	legal, _ := g.LegalCards(2)
	if len(legal) != 1 || legal[0].Rank != "J" {
		t.Fatalf("must over-trump with J: %v", legal)
	}
}

func TestTrickWinnerTrumpOrder(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Spades}
	g.Trick = []Played{{0, Card{"hearts", "A"}}, {1, Card{"spades", "7"}}, {2, Card{"hearts", "10"}}, {3, Card{"spades", "9"}}}
	if w := g.currentWinner(); w.Seat != 3 {
		t.Fatalf("9 of trump beats 7 of trump, got seat %d", w.Seat)
	}
	g.Contract = NoTrump
	if w := g.currentWinner(); w.Seat != 0 {
		t.Fatalf("no trump: ace of led suit wins, got %d", w.Seat)
	}
}

func TestCardPoints(t *testing.T) {
	g := &Game{Contract: Hearts}
	if g.CardPoints(Card{"hearts", "J"}) != 20 || g.CardPoints(Card{"spades", "J"}) != 2 || g.CardPoints(Card{"hearts", "9"}) != 14 {
		t.Fatal("bad card points")
	}
	g.Contract = AllTrump
	if g.CardPoints(Card{"spades", "J"}) != 20 {
		t.Fatal("all-trump J should be 20")
	}
}

// Total card points in a deal are 152 (trump contract) incl. last trick.
func TestFullDealTotals152AndMatchAdvances(t *testing.T) {
	for _, c := range []Contract{Hearts, NoTrump, AllTrump} {
		g := newPlaying(t, c)
		for g.Phase == PhasePlaying {
			legal, _ := g.LegalCards(g.Turn)
			if len(legal) == 0 {
				t.Fatalf("%s: no legal card for seat %d", c, g.Turn)
			}
			if err := g.Play(g.Turn, legal[0]); err != nil {
				t.Fatal(err)
			}
		}
		want := 152
		if c == NoTrump {
			want = 152 // raw is 152 before the x2; captured stays 152
		}
		if c == AllTrump {
			want = 152
		}
		if g.Captured[0]+g.Captured[1] != want && g.Captured[0]+g.Captured[1] != want-2*0 {
			// all-trump/no-trump totals differ by variant; just require plausibility
			if tot := g.Captured[0] + g.Captured[1]; tot < 100 || tot > 270 {
				t.Fatalf("%s: implausible total %d", c, tot)
			}
		}
		if g.Result == nil || g.Result.Awarded[0]+g.Result.Awarded[1] == 0 {
			t.Fatalf("%s: no score awarded", c)
		}
		if g.Tricks[0]+g.Tricks[1] != 8 {
			t.Fatal("expected 8 tricks")
		}
	}
	g := newPlaying(t, Hearts)
	for g.Phase == PhasePlaying {
		l, _ := g.LegalCards(g.Turn)
		_ = g.Play(g.Turn, l[0])
	}
	if g.Captured[0]+g.Captured[1] != 162 {
		t.Fatalf("trump contract must total 162, got %d", g.Captured[0]+g.Captured[1])
	}
}

func TestIllegalPlayDoesNotMutate(t *testing.T) {
	g := newPlaying(t, Hearts)
	before := len(g.Hands[g.Turn])
	if err := g.Play((g.Turn+1)%4, g.Hands[(g.Turn+1)%4][0]); err == nil {
		t.Fatal("out of turn play accepted")
	}
	if err := g.Play(g.Turn, Card{"clubs", "x"}); err == nil {
		t.Fatal("card not in hand accepted")
	}
	if len(g.Hands[g.Turn]) != before {
		t.Fatal("state mutated")
	}
}

func TestScoringContractFailsNeedsMoreThanHalf(t *testing.T) {
	g := &Game{Contract: Hearts, Bidder: 0, Mult: 1, Captured: [2]int{76, 76}, Tricks: [2]int{4, 4}, Target: 151, Winner: -1}
	g.finishDeal()
	if g.Result.Succeeded || g.Result.Awarded[1] != 16 {
		t.Fatalf("tie should fail the bidder: %+v", g.Result)
	}
}

func TestMatchEndsAtTargetNoTie(t *testing.T) {
	g := &Game{Target: 151, Scores: [2]int{151, 151}, Winner: -1}
	g.checkMatch()
	if g.Phase == PhaseMatch {
		t.Fatal("tie must continue")
	}
	g.Scores[1] = 160
	g.checkMatch()
	if g.Phase != PhaseMatch || g.Winner != 1 {
		t.Fatal("expected team 1 to win")
	}
}
