package belot

import (
	"encoding/json"
	"reflect"
	"testing"
)

func declarationHand() []Card {
	return []Card{{"clubs", "J"}, {"diamonds", "J"}, {"hearts", "J"}, {"spades", "J"}, {"clubs", "7"}, {"clubs", "8"}, {"clubs", "9"}, {"clubs", "10"}}
}
func declarationID(t *testing.T, g *Game, seat int, kind string, points, length int) string {
	t.Helper()
	for _, d := range g.DeclarationOptions(seat) {
		if d.Kind == kind && d.Points == points && len(d.Cards) == length {
			return d.ID
		}
	}
	t.Fatalf("no declaration %s %d/%d", kind, points, length)
	return ""
}
func TestDeclarationChoicesValidateOverlapAndTimingAtomically(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Turn: 0}
	g.Hands[0] = declarationHand()
	four := declarationID(t, g, 0, "four", 200, 4)
	five := declarationID(t, g, 0, "sequence", 100, 5)
	fourSeq := declarationID(t, g, 0, "sequence", 50, 4)
	before, _ := json.Marshal(g)
	if err := g.PlayDeclared(0, g.Hands[0][0], []string{four, five}, false); err == nil {
		t.Fatal("overlapping four and sequence accepted")
	}
	after, _ := json.Marshal(g)
	if string(before) != string(after) {
		t.Fatal("rejected declarations mutated state")
	}
	if err := g.PlayDeclared(0, g.Hands[0][0], []string{four, fourSeq}, false); err != nil {
		t.Fatal(err)
	}
	if len(g.PublicDeclarations()) != 2 || g.DeclarationPoints() != ([2]int{250, 0}) {
		t.Fatalf("valid nonoverlap choices not counted: %v", g.PublicDeclarations())
	}
	if len(g.DeclarationOptions(0)) != 0 {
		t.Fatal("late declaration options exposed")
	}
	public, _ := json.Marshal(g.PublicDeclarations())
	if string(public) == "" {
		t.Fatal("missing public summary")
	}
	var summaries []map[string]any
	_ = json.Unmarshal(public, &summaries)
	for _, s := range summaries {
		if _, ok := s["cards"]; ok {
			t.Fatal("public summary leaks cards")
		}
		if _, ok := s["id"]; ok {
			t.Fatal("public summary leaks combination id")
		}
	}
}

func TestSequenceAndFourCategoriesResolveIndependently(t *testing.T) {
	g := &Game{Contract: AllTrump}
	g.declared = []declaredCombination{
		{Seat: 0, Declaration: Declaration{Kind: "sequence", Points: 20, length: 3, high: 5}},
		{Seat: 1, Declaration: Declaration{Kind: "sequence", Points: 50, length: 4, high: 4}},
		{Seat: 2, Declaration: Declaration{Kind: "four", Points: 200, high: 7}},
		{Seat: 3, Declaration: Declaration{Kind: "four", Points: 150, high: 6}},
	}
	if got := g.DeclarationPoints(); got != ([2]int{200, 50}) {
		t.Fatalf("independent categories: %v", got)
	}
	g.declared = append(g.declared, declaredCombination{Seat: 0, Declaration: Declaration{Kind: "sequence", Points: 50, length: 4, high: 4}})
	if got := g.DeclarationPoints(); got != ([2]int{200, 0}) {
		t.Fatalf("equal strongest sequences must cancel: %v", got)
	}
}

func TestBelotRequiresEligibleFirstKingOrQueen(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Turn: 0}
	g.Hands[0] = []Card{{"hearts", "K"}, {"hearts", "Q"}, {"clubs", "7"}}
	if !g.CanAnnounceBelot(0, Card{"hearts", "K"}) {
		t.Fatal("eligible belot unavailable")
	}
	if err := g.PlayDeclared(0, Card{"hearts", "K"}, nil, true); err != nil {
		t.Fatal(err)
	}
	if g.BelotPoints != ([2]int{20, 0}) {
		t.Fatal("belot not awarded")
	}
	g = &Game{Phase: PhasePlaying, Contract: AllTrump, Turn: 0, Trick: []Played{{3, Card{"diamonds", "7"}}}}
	g.Hands[0] = []Card{{"hearts", "K"}, {"hearts", "Q"}, {"clubs", "7"}}
	before := append([]Card{}, g.Hands[0]...)
	if g.CanAnnounceBelot(0, Card{"hearts", "K"}) {
		t.Fatal("off-suit all-trump discard allowed belot")
	}
	if err := g.PlayDeclared(0, Card{"hearts", "K"}, nil, true); err == nil {
		t.Fatal("invalid belot accepted")
	}
	if !reflect.DeepEqual(before, g.Hands[0]) || g.BelotPoints != ([2]int{}) {
		t.Fatal("invalid belot mutated hand/points")
	}
}

func TestNoTrumpRejectsAnnouncementsWithoutMutatingCard(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: NoTrump, Turn: 0}
	g.Hands[0] = declarationHand()
	if len(g.DeclarationOptions(0)) != 0 {
		t.Fatal("no-trump offered declarations")
	}
	if err := g.PlayDeclared(0, g.Hands[0][0], []string{"four:J"}, false); err == nil {
		t.Fatal("no-trump declaration accepted")
	}
	if len(g.Hands[0]) != 8 || len(g.Trick) != 0 {
		t.Fatal("rejected no-trump declaration played card")
	}
}

func TestDeclarationOptionsOnlyOfferMaximalRunsAndFourAlternatives(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Turn: 0}
	g.Hands[0] = declarationHand()
	if options := g.DeclarationOptions(0); len(options) != 3 {
		t.Fatalf("want four J, full quint, and four-compatible quart, got %v", options)
	}
	shorter := declarationID(t, g, 0, "sequence", 50, 4)
	if err := g.PlayDeclared(0, g.Hands[0][0], []string{shorter}, false); err == nil {
		t.Fatal("truncated run accepted without choosing the four that requires truncation")
	}
	g.Hands[0] = NewDeck()[:8]
	if options := g.DeclarationOptions(0); len(options) != 1 || len(options[0].Cards) != 8 {
		t.Fatalf("one long suit must not split into nested runs: %v", options)
	}
}

func TestDefaultPlayDoesNotAutoDeclareAndLateSelectionDoesNotPlayCard(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Turn: 0}
	g.Hands[0] = declarationHand()
	id := declarationID(t, g, 0, "four", 200, 4)
	if err := g.Play(0, g.Hands[0][0]); err != nil {
		t.Fatal(err)
	}
	if len(g.PublicDeclarations()) != 0 {
		t.Fatal("default play silently announced")
	}
	g.Turn = 0
	g.Trick = nil
	before := append([]Card{}, g.Hands[0]...)
	if err := g.PlayDeclared(0, g.Hands[0][0], []string{id}, false); err == nil {
		t.Fatal("late declaration accepted")
	}
	if !reflect.DeepEqual(g.Hands[0], before) || len(g.PublicDeclarations()) != 0 || len(g.Trick) != 0 {
		t.Fatal("late declaration mutated game")
	}
}

func TestBelotCannotBeClaimedAfterSkippingFirstPairCard(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: AllTrump, Turn: 0}
	g.Hands[0] = []Card{{"hearts", "K"}, {"hearts", "Q"}}
	if err := g.Play(0, Card{"hearts", "K"}); err != nil {
		t.Fatal(err)
	}
	g.Turn = 0
	g.Trick = nil
	if g.CanAnnounceBelot(0, Card{"hearts", "Q"}) {
		t.Fatal("second pair card offered belot after missing first")
	}
	if err := g.PlayDeclared(0, Card{"hearts", "Q"}, nil, true); err == nil {
		t.Fatal("belot accepted on second pair card")
	}
	if len(g.Hands[0]) != 1 {
		t.Fatal("invalid belot consumed second card")
	}
}

func TestInvalidCardWithValidDeclarationsLeavesAnnouncementStateUntouched(t *testing.T) {
	g := &Game{Phase: PhasePlaying, Contract: Hearts, Turn: 0}
	g.Hands[0] = declarationHand()
	id := declarationID(t, g, 0, "four", 200, 4)
	if err := g.PlayDeclared(0, Card{"spades", "A"}, []string{id}, false); err == nil {
		t.Fatal("unheld card accepted")
	}
	if len(g.Hands[0]) != 8 || len(g.PublicDeclarations()) != 0 || g.BelotPoints != ([2]int{}) {
		t.Fatal("invalid card committed declarations")
	}
}
