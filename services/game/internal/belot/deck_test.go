package belot

import (
	"bytes"
	"testing"
)

func TestDeckContains32UniqueCards(t *testing.T) {
	if err := ValidateDeck(NewDeck()); err != nil {
		t.Fatal(err)
	}
}

func TestShuffleAndDeal(t *testing.T) {
	// A deterministic source tests the algorithm without asserting random order.
	src := bytes.NewReader(bytes.Repeat([]byte{0}, 512))
	deck, err := Shuffle(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeck(deck); err != nil {
		t.Fatal(err)
	}
	hands, err := Deal(deck)
	if err != nil {
		t.Fatal(err)
	}
	used := make(map[Card]bool)
	for _, hand := range hands {
		if len(hand) != 8 {
			t.Fatalf("wanted 8 cards, got %d", len(hand))
		}
		for _, card := range hand {
			if used[card] {
				t.Fatalf("duplicate card in hands: %v", card)
			}
			used[card] = true
		}
	}
	if len(used) != 32 {
		t.Fatalf("used %d cards", len(used))
	}
}

func TestInvalidDeckRejected(t *testing.T) {
	deck := NewDeck()
	deck[0] = deck[1]
	if _, err := Deal(deck); err == nil {
		t.Fatal("expected duplicate card error")
	}
	if err := ValidateDeck(NewDeck()[:31]); err == nil {
		t.Fatal("expected short deck error")
	}
}
