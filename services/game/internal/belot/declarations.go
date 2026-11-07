package belot

import "fmt"

// Declaration options are private to their owner's seat, never a public game projection.
type Declaration struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Points int    `json:"points"`
	Cards  []Card `json:"cards"`
	length int
	high   int
}

type DeclaredSummary struct {
	Seat   int    `json:"seat"`
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type declaredCombination struct {
	Seat        int
	Declaration Declaration
}

// DeclarationOptions offers maximal runs and the shorter runs made available by
// choosing an overlapping four. It does not offer every nested sub-sequence.
func (g *Game) DeclarationOptions(seat int) []Declaration {
	if seat < 0 || seat > 3 || g.Phase != PhasePlaying || g.Contract == NoTrump || len(g.Hands[seat]) != 8 {
		return nil
	}
	hand := g.Hands[seat]
	present := map[Card]bool{}
	for _, c := range hand {
		present[c] = true
	}
	var options []Declaration
	var fourRanks []Rank
	for _, rank := range Ranks {
		var cards []Card
		for _, c := range hand {
			if c.Rank == rank {
				cards = append(cards, c)
			}
		}
		if len(cards) != 4 || rank == "7" || rank == "8" {
			continue
		}
		points := 100
		if rank == "9" {
			points = 150
		}
		if rank == "J" {
			points = 200
		}
		options = append(options, Declaration{ID: "four:" + string(rank), Kind: "four", Name: "Каре", Points: points, Cards: cards, high: rankIndex(rank, trumpOrder)})
		fourRanks = append(fourRanks, rank)
	}
	seen := map[string]bool{}
	for mask := 0; mask < (1 << len(fourRanks)); mask++ {
		excluded := map[Rank]bool{}
		for i, r := range fourRanks {
			if mask&(1<<i) != 0 {
				excluded[r] = true
			}
		}
		for _, suit := range Suits {
			start := 0
			for start < len(Ranks) {
				if excluded[Ranks[start]] || !present[Card{suit, Ranks[start]}] {
					start++
					continue
				}
				end := start
				var cards []Card
				for end < len(Ranks) && !excluded[Ranks[end]] && present[Card{suit, Ranks[end]}] {
					cards = append(cards, Card{suit, Ranks[end]})
					end++
				}
				if len(cards) >= 3 {
					id := fmt.Sprintf("sequence:%s:%d:%d", suit, start, len(cards))
					if !seen[id] {
						points, name := 20, "Терца"
						if len(cards) == 4 {
							points, name = 50, "Кварта"
						}
						if len(cards) >= 5 {
							points, name = 100, "Квинта"
						}
						options = append(options, Declaration{ID: id, Kind: "sequence", Name: name, Points: points, Cards: cards, length: len(cards), high: end - 1})
						seen[id] = true
					}
				}
				start = end
			}
		}
	}
	return options
}

func (g *Game) selectedDeclarations(seat int, ids []string) ([]Declaration, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	options := g.DeclarationOptions(seat)
	selected := make([]Declaration, 0, len(ids))
	used := map[Card]bool{}
	for _, id := range ids {
		var found *Declaration
		for i := range options {
			if options[i].ID == id {
				found = &options[i]
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("Анонсът не е наличен или вече е късно за обявяване.")
		}
		for _, c := range found.Cards {
			if used[c] {
				return nil, fmt.Errorf("Една карта не може да участва в два избрани анонса.")
			}
			used[c] = true
		}
		selected = append(selected, *found)
	}
	excluded := map[Card]bool{}
	for _, d := range selected {
		if d.Kind == "four" {
			for _, c := range d.Cards {
				excluded[c] = true
			}
		}
	}
	for _, d := range selected {
		if d.Kind != "sequence" {
			continue
		}
		first, last := rankIndex(d.Cards[0].Rank, Ranks[:]), rankIndex(d.Cards[len(d.Cards)-1].Rank, Ranks[:])
		for _, neighbor := range []int{first - 1, last + 1} {
			if neighbor < 0 || neighbor >= len(Ranks) {
				continue
			}
			card := Card{d.Cards[0].Suit, Ranks[neighbor]}
			for _, held := range g.Hands[seat] {
				if held == card && !excluded[card] {
					return nil, fmt.Errorf("Обявете цялата поредица или изберете карето, което я прекъсва.")
				}
			}
		}
	}
	return selected, nil
}

func (g *Game) PublicDeclarations() []DeclaredSummary {
	summaries := make([]DeclaredSummary, 0, len(g.declared))
	for _, d := range g.declared {
		summaries = append(summaries, DeclaredSummary{Seat: d.Seat, Kind: d.Declaration.Kind, Name: d.Declaration.Name, Points: d.Declaration.Points})
	}
	return summaries
}

// DeclarationPoints resolves sequences and fours independently. Matching strongest
// sequences cancel the entire sequence category, including lesser declarations.
func (g *Game) DeclarationPoints() [2]int {
	var points [2]int
	for _, kind := range []string{"sequence", "four"} {
		var best [2]*Declaration
		var sums [2]int
		for _, d := range g.declared {
			c := d.Declaration
			if c.Kind != kind {
				continue
			}
			team := Team(d.Seat)
			sums[team] += c.Points
			if best[team] == nil || declarationCompare(c, *best[team]) > 0 {
				copy := c
				best[team] = &copy
			}
		}
		switch {
		case best[0] == nil && best[1] == nil:
		case best[1] == nil:
			points[0] += sums[0]
		case best[0] == nil:
			points[1] += sums[1]
		default:
			if cmp := declarationCompare(*best[0], *best[1]); cmp > 0 {
				points[0] += sums[0]
			} else if cmp < 0 {
				points[1] += sums[1]
			}
		}
	}
	return points
}

func declarationCompare(a, b Declaration) int {
	if a.Kind == "sequence" && a.length != b.length {
		return a.length - b.length
	}
	return a.high - b.high
}

func (g *Game) CanAnnounceBelot(seat int, card Card) bool {
	if seat < 0 || seat > 3 || g.Phase != PhasePlaying || g.Turn != seat || g.Contract == NoTrump || (card.Rank != "K" && card.Rank != "Q") {
		return false
	}
	trump, oneTrump := g.Contract.trumpSuit()
	if g.Contract != AllTrump && (!oneTrump || card.Suit != trump) {
		return false
	}
	haveK, haveQ := false, false
	for _, held := range g.Hands[seat] {
		if held.Suit == card.Suit {
			haveK = haveK || held.Rank == "K"
			haveQ = haveQ || held.Rank == "Q"
		}
	}
	if !haveK || !haveQ {
		return false
	}
	if len(g.Trick) > 0 && card.Suit != g.Trick[0].Card.Suit && !(oneTrump && card.Suit == trump) {
		return false
	}
	legal, _ := g.LegalCards(seat)
	for _, c := range legal {
		if c == card {
			return true
		}
	}
	return false
}
