package belot

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
)

type Suit string
type Rank string

type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

var Suits = [...]Suit{"clubs", "diamonds", "hearts", "spades"}
var Ranks = [...]Rank{"7", "8", "9", "10", "J", "Q", "K", "A"}

// NewDeck creates the 32-card Bulgarian Belot deck. It does not expose private hands.
func NewDeck() []Card {
	deck := make([]Card, 0, 32)
	for _, suit := range Suits {
		for _, rank := range Ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

// Shuffle uses cryptographic randomness. No client chooses or supplies the deal.
func Shuffle(src io.Reader) ([]Card, error) {
	if src == nil {
		src = rand.Reader
	}
	deck := NewDeck()
	for i := len(deck) - 1; i > 0; i-- {
		j, err := rand.Int(src, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, fmt.Errorf("shuffle: %w", err)
		}
		k := int(j.Int64())
		deck[i], deck[k] = deck[k], deck[i]
	}
	return deck, nil
}

// Deal distributes 8 cards to each seat. This does not yet implement the
// bidding-stage 5+3 card delivery, which is a subsequent engine milestone.
func Deal(deck []Card) ([4][]Card, error) {
	var hands [4][]Card
	if err := ValidateDeck(deck); err != nil {
		return hands, err
	}
	for i, card := range deck {
		hands[i%4] = append(hands[i%4], card)
	}
	return hands, nil
}

// ValidateDeck rejects duplicates, unsupported ranks/suits and incomplete decks.
func ValidateDeck(deck []Card) error {
	if len(deck) != 32 {
		return fmt.Errorf("expected 32 cards, got %d", len(deck))
	}
	expected := make(map[Card]bool, 32)
	for _, card := range NewDeck() {
		expected[card] = true
	}
	seen := make(map[Card]bool, 32)
	for _, card := range deck {
		if !expected[card] {
			return fmt.Errorf("invalid card: %v", card)
		}
		if seen[card] {
			return fmt.Errorf("duplicate card: %v", card)
		}
		seen[card] = true
	}
	if len(seen) != len(expected) {
		return errors.New("incomplete deck")
	}
	return nil
}
