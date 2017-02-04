package internal

import (
	"math/rand"
	"time"
)

//FrenchCard is a struct for cards
type FrenchCard struct {
	value string
	suit  string
}

func NewCard(v, s string) *FrenchCard {
	return &FrenchCard{
		value: v,
		suit:  s,
	}
}

func (c *FrenchCard) Value() string {
	return c.value
}
func (c *FrenchCard) Suit() string {
	return c.suit
}

type FrenchDeck struct {
	deck   []Card
	cursor int
	r      *rand.Rand
}

func NewDeck() *FrenchDeck {
	self := &FrenchDeck{
		cursor: 0,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, v := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"} {
		for _, s := range []string{"\u2665", "\u2663", "\u2666", "\u2660"} {
			self.deck = append(self.deck, NewCard(v, s))
		}
	}

	return self
}

func NewHand() *FrenchDeck {
	return &FrenchDeck{
		cursor: 0,
	}
}

// Deck generates, shuffles, and deals a deck of 52 cards
func (self *FrenchDeck) Shuffle() {
	// Shuffle the deck using Fisher-Yates
	// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := 1; i < len(self.deck); i++ {
		n := self.r.Intn(i)
		self.deck[n], self.deck[i] = self.deck[i], self.deck[n]
	}
}

func (self *FrenchDeck) Insert(card Card) {
	self.deck = append(self.deck, card)
}

func (self *FrenchDeck) Draw() Card {
	card := self.deck[0]
	self.deck = self.deck[1:]
	return card
}

func (self *FrenchDeck) Next() Card {
	if self.cursor == len(self.deck) {
		self.cursor = 0
		return nil
	}
	card := self.deck[self.cursor]
	self.cursor = (self.cursor + 1) % (len(self.deck) + 1)
	return card
}

func (self *FrenchDeck) Peek() Card {
	if self.cursor == len(self.deck) {
		return nil
	}
	return self.deck[0]
}

func (self *FrenchDeck) Size() int {
	return len(self.deck)
}
