package belot

import (
	"math/rand"
	"reflect"
	"testing"
)

// Golden rules: https://belot.bg/belot/rules/ (bidding, scoring, hanging, capot).
// Pair rounding examples: https://belotarena.com/kalkulator-tochki-belot.
func TestRedoubleKeepsAuctionOpenForThreePasses(t *testing.T) {
	g := NewGame(3, 151)
	if err := g.StartDealWith(NewDeck()); err != nil {
		t.Fatal(err)
	}
	for _, b := range []struct {
		seat   int
		action string
		c      Contract
	}{{0, "bid", Clubs}, {1, "double", ""}, {2, "redouble", ""}} {
		if err := g.Bid(b.seat, b.action, b.c); err != nil {
			t.Fatal(err)
		}
	}
	if g.Phase != PhaseBidding || g.Turn != 3 {
		t.Fatalf("recontra closed auction: phase=%s turn=%d", g.Phase, g.Turn)
	}
	if err := g.Bid(3, "bid", Hearts); err != nil {
		t.Fatal("higher contract must remain possible:", err)
	}
	if g.Mult != 1 {
		t.Fatal("higher contract must reset multiplier")
	}
	for _, s := range []int{0, 1, 2} {
		if err := g.Bid(s, "pass", ""); err != nil {
			t.Fatal(err)
		}
	}
	if g.Phase != PhasePlaying {
		t.Fatal("auction did not close after 3 passes")
	}
}

func TestDealUsesThreeTwoThreePacketsFromDealersRight(t *testing.T) {
	deck := NewDeck()
	g := NewGame(1, 151)
	if err := g.StartDealWith(deck); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		s := (2 + i) % 4
		want := append(append([]Card{}, deck[i*3:i*3+3]...), deck[12+i*2:14+i*2]...)
		if !sameCards(g.Hands[s], want) {
			t.Fatalf("seat %d initial packet=%v want=%v", s, g.Hands[s], want)
		}
	}
	if err := g.Bid(2, "bid", Hearts); err != nil {
		t.Fatal(err)
	}
	for _, s := range []int{3, 0, 1} {
		if err := g.Bid(s, "pass", ""); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 4; i++ {
		s := (2 + i) % 4
		for _, c := range deck[20+i*3 : 23+i*3] {
			if !hasCard(g.Hands[s], c) {
				t.Fatalf("seat %d missing final packet card %v", s, c)
			}
		}
	}
}

func hasCard(cards []Card, c Card) bool {
	for _, x := range cards {
		if x == c {
			return true
		}
	}
	return false
}
func sameCards(a, b []Card) bool {
	if len(a) != len(b) {
		return false
	}
	for _, c := range a {
		if !hasCard(b, c) {
			return false
		}
	}
	return true
}

func TestContractSpecificPairRoundingAndWholePot(t *testing.T) {
	for _, tt := range []struct {
		name string
		c    Contract
		raw  [2]int
		mult int
		want [2]int
	}{
		{"trump lower six rounds up", Hearts, [2]int{96, 66}, 1, [2]int{9, 7}},
		{"trump lower six reversed", Hearts, [2]int{86, 76}, 1, [2]int{8, 8}},
		{"all-trump lower four rounds up", AllTrump, [2]int{134, 124}, 1, [2]int{13, 13}},
		{"no-trump doubles before rounding", NoTrump, [2]int{79, 51}, 1, [2]int{16, 10}},
		{"failed contract rounds whole pot", AllTrump, [2]int{124, 134}, 1, [2]int{0, 26}},
		{"contra rounds whole pot", AllTrump, [2]int{134, 124}, 2, [2]int{52, 0}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{Contract: tt.c, Bidder: 0, Mult: tt.mult, Captured: tt.raw, Tricks: [2]int{4, 4}, Target: 151, Winner: -1}
			g.finishDeal()
			if g.Result.Awarded != tt.want {
				t.Fatalf("got %v want %v", g.Result.Awarded, tt.want)
			}
		})
	}
}

func TestCapotBonusNotDoubledByNoTrumpAndCannotEndMatch(t *testing.T) {
	for _, tt := range []struct {
		c               Contract
		captured, award int
	}{{Hearts, 162, 25}, {NoTrump, 130, 35}, {AllTrump, 258, 35}} {
		g := &Game{Contract: tt.c, Bidder: 0, Mult: 1, Captured: [2]int{tt.captured, 0}, Tricks: [2]int{8, 0}, Scores: [2]int{150, 0}, Target: 151, Winner: -1}
		g.finishDeal()
		if g.Result.Awarded != ([2]int{tt.award, 0}) {
			t.Fatalf("%s capot score=%v", tt.c, g.Result.Awarded)
		}
		if g.Phase != PhaseDealOver {
			t.Fatalf("%s cannot finish match on capot", tt.c)
		}
	}
}

func TestRejectedBidLeavesEntireGameUnchanged(t *testing.T) {
	g := NewGame(3, 151)
	_ = g.StartDealWith(NewDeck())
	_ = g.Bid(0, "bid", Hearts)
	before := *g
	before.Hands = g.Hands
	before.Log = append([]string{}, g.Log...)
	if err := g.Bid(1, "bid", Clubs); err == nil {
		t.Fatal("lower contract accepted")
	}
	if !reflect.DeepEqual(before, *g) {
		t.Fatal("rejected bid mutated game")
	}
}

func TestHangingPointsCarryAndAccumulateUntilAnUnequalDeal(t *testing.T) {
	g := &Game{Contract: Hearts, Bidder: 0, Mult: 1, Captured: [2]int{81, 81}, Tricks: [2]int{4, 4}, Target: 151, Winner: -1}
	g.finishDeal()
	if !g.Result.Hanging || g.Result.Awarded != ([2]int{0, 8}) || g.Pending != 8 {
		t.Fatalf("first hanging deal: %+v pending=%d", g.Result, g.Pending)
	}
	g.Mult = 2
	g.Captured = [2]int{81, 81}
	g.finishDeal()
	if !g.Result.Hanging || g.Result.Awarded != ([2]int{}) || g.Pending != 40 {
		t.Fatalf("doubled hanging deal: %+v pending=%d", g.Result, g.Pending)
	}
	g.Mult = 1
	g.Captured = [2]int{96, 66}
	g.finishDeal()
	if g.Result.Hanging || g.Result.Awarded != ([2]int{49, 7}) || g.Pending != 0 || g.Result.CarryAwarded != 40 {
		t.Fatalf("resolved carry: %+v pending=%d", g.Result, g.Pending)
	}
}

func TestHangingRoundingExamplesFromPublishedRules(t *testing.T) {
	for _, tt := range []struct {
		c     Contract
		raw   int
		want  [2]int
		carry int
	}{{Hearts, 106, [2]int{0, 10}, 11}, {AllTrump, 154, [2]int{0, 15}, 16}, {NoTrump, 65, [2]int{0, 13}, 13}} {
		g := &Game{Contract: tt.c, Bidder: 0, Mult: 1, Captured: [2]int{tt.raw, tt.raw}, Tricks: [2]int{4, 4}, Target: 151, Winner: -1}
		if tt.c != NoTrump {
			g.Captured[0] -= 50
			g.declared = []declaredCombination{{Seat: 0, Declaration: Declaration{Kind: "sequence", Name: "Кварта", Points: 50, length: 4, high: 4}}}
		}
		g.finishDeal()
		if g.Result.Awarded != tt.want || g.Pending != tt.carry {
			t.Fatalf("%s hanging got %v/%d want %v/%d", tt.c, g.Result.Awarded, g.Pending, tt.want, tt.carry)
		}
	}
}

func TestTenThousandSeededDealsPreserveCardsTurnsAndPointTotals(t *testing.T) {
	rng := rand.New(rand.NewSource(303))
	contracts := []Contract{Clubs, Diamonds, Hearts, Spades, NoTrump, AllTrump}
	for n := 0; n < 10000; n++ {
		c := contracts[n%len(contracts)]
		g := NewGame(n%4, 151)
		if err := g.StartDeal(rng); err != nil {
			t.Fatal(err)
		}
		if err := g.Bid(g.Turn, "bid", c); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 3; i++ {
			if err := g.Bid(g.Turn, "pass", ""); err != nil {
				t.Fatal(err)
			}
		}
		seen := map[Card]bool{}
		for i := 0; i < 32; i++ {
			if g.Phase != PhasePlaying {
				t.Fatalf("deal %d ended early at card %d", n, i)
			}
			seat := g.Turn
			legal, _ := g.LegalCards(seat)
			if len(legal) == 0 {
				t.Fatal("empty legal action set")
			}
			card := legal[rng.Intn(len(legal))]
			if seen[card] {
				t.Fatalf("deal %d duplicated card %v", n, card)
			}
			seen[card] = true
			if err := g.Play(seat, card); err != nil {
				t.Fatal(err)
			}
		}
		total := 162
		if c == NoTrump {
			total = 130
		}
		if c == AllTrump {
			total = 258
		}
		if len(seen) != 32 || g.Tricks[0]+g.Tricks[1] != 8 || g.Captured[0]+g.Captured[1] != total {
			t.Fatalf("deal %d invariant failure", n)
		}
		if g.Phase != PhaseDealOver && g.Phase != PhaseMatch {
			t.Fatalf("deal %d never ended", n)
		}
	}
}

func TestPublicWinnerAndComparator(t *testing.T) {
	g := &Game{Contract: Hearts}
	if g.WinningSeat() != -1 {
		t.Fatal("empty trick has winner")
	}
	g.Trick = []Played{{0, Card{"clubs", "A"}}, {1, Card{"hearts", "7"}}}
	if g.WinningSeat() != 1 || !g.Beats("clubs", Card{"hearts", "7"}, Card{"hearts", "9"}) {
		t.Fatal("public comparator disagrees with trump order")
	}
}

func TestScoreLedgerComponentsSumToRawPoints(t *testing.T) {
	for _, tt := range []struct {
		c        Contract
		captured [2]int
		tricks   [2]int
	}{{Hearts, [2]int{96, 66}, [2]int{5, 3}}, {NoTrump, [2]int{130, 0}, [2]int{8, 0}}, {AllTrump, [2]int{258, 0}, [2]int{8, 0}}} {
		g := &Game{Contract: tt.c, Bidder: 0, Mult: 1, Captured: tt.captured, Tricks: tt.tricks, BelotPoints: [2]int{20, 0}, Target: 151, Winner: -1}
		if tt.c == NoTrump {
			g.BelotPoints = [2]int{}
		}
		g.finishDeal()
		r := g.Result
		for team := 0; team < 2; team++ {
			if r.RawPoints[team] != r.CardPoints[team]+r.DeclarationPoints[team]+r.BelotPoints[team]+r.CapotPoints[team] {
				t.Fatalf("%s components do not sum: %+v", tt.c, r)
			}
		}
	}
}

func TestAllPassPreservesAndReportsPendingPoints(t *testing.T) {
	g := NewGame(3, 151)
	g.Pending = 11
	if err := g.StartDealWith(NewDeck()); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		if err := g.Bid(g.Turn, "pass", ""); err != nil {
			t.Fatal(err)
		}
	}
	if g.Pending != 11 || g.Result.Carried != 11 {
		t.Fatalf("all-pass dropped pending points: %d / %+v", g.Pending, g.Result)
	}
	if err := g.NextDeal(nil); err != nil {
		t.Fatal(err)
	}
	if g.Pending != 11 {
		t.Fatal("next deal dropped pending points")
	}
}
