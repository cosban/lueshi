package games

import (
	"testing"

	"github.com/cosban/assert"
	"github.com/cosban/lueshi/internal"
)

func TestNewInstance(t *testing.T) {
	blkjk := NewBLKJK()
	assert.Equals(t, "Blackjack", blkjk.Name())
}

func TestNewPlayer(t *testing.T) {
	assert := assert.New(t)
	player := NewPlayer("1")

	assert.Equals("1", player.ID)
	assert.Equals(nil, player.hand.Peek())
	assert.Equals(100, player.bank)
}

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	player := NewPlayer("1")

	b := NewBLKJK()
	assert.True(b.addPlayer(player))
	assert.True(b.playerExists(player.ID))
	assert.Equals(player, b.getPlayer(player.ID))
	assert.Equals(1, len(b.players))
}

func TestRemovePlayer(t *testing.T) {
	assert := assert.New(t)
	player := NewPlayer("1")

	b := NewBLKJK()
	assert.False(b.removePlayer(player))
	assert.True(b.addPlayer(player))
	assert.True(b.removePlayer(player))
	assert.False(b.playerExists(player.ID))
}

func TestDealhand(t *testing.T) {
	assert := assert.New(t)
	player := NewPlayer("1")

	b := NewBLKJK()
	b.addPlayer(player)
	b.dealNewHand()
	assert.Equals(2, b.dealer.Size())
	assert.Equals(2, player.hand.Size())
}

func TestGetHand(t *testing.T) {
	assert := assert.New(t)

	deck := internal.NewHand()
	assert.Equals("", getHand(deck))
	deck.Insert(internal.NewCard("K", "D"))
	assert.Equals("KD", getHand(deck))
	deck.Insert(internal.NewCard("Q", "C"))
	assert.Equals("KD, QC", getHand(deck))
}

func TestHandTotal(t *testing.T) {
	assert := assert.New(t)
	deck := internal.NewHand()

	deck.Insert(internal.NewCard("A", "Z"))
	assert.Equals(11, handTotal(deck))
	deck.Insert(internal.NewCard("A", "Y"))
	assert.Equals(12, handTotal(deck))
	deck.Insert(internal.NewCard("K", "X"))
	assert.Equals(12, handTotal(deck))
}
