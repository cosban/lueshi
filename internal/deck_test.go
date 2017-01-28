package internal

import "testing"

func TestNewDeck(t *testing.T) {
	deck := NewDeck()
	if deck.Size() != 52 {
		t.Fatal("The deck is the wrong size! expected %d received %d", 52, deck.Size())
	}
}

func TestNewHand(t *testing.T) {
	deck := NewHand()
	if deck.Size() != 0 {
		t.Fatalf("The deck should have no cards in it!")
	}
	if deck.Peek() != nil {
		t.Fatalf("The deck should have no cards in it!")
	}
}

func TestShuffle(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle()
	if deck.Size() != 52 {
		t.Fatalf("The deck is the wrong size! expected: %d received: %d", 52, deck.Size())
	}
}

func TestPeek(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle()
	c := deck.Peek()
	if c != deck.Peek() {
		t.Fatal("Peek is nondeterministic!")
	}
}

func TestDraw(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle()
	c := deck.Draw()

	if c == deck.Peek() {
		t.Fatal("Failed to remove card from the deck during a draw")
	}
}

func TestInsert(t *testing.T) {
	deck := NewHand()
	if deck.Peek() != nil {
		t.Fatalf("The deck should have no cards in it!")
	}
	c := &FrenchCard{"FAKE", "CARD"}
	deck.Insert(c)
	if deck.Peek() != c {
		t.Fatalf("Failed to insert the card!")
	}
}
