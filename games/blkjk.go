package games

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"strconv"

	"github.com/cosban/lueshi/command"
	"github.com/cosban/lueshi/internal"
)

//BLKJK is shorthand for Blackjack
type BLKJK struct {
	name        string
	description string
	rules       string
	players     []*CardPlayer
	dealer      internal.Deck
	deck        internal.Deck
	cursor      int

	actions    map[string]command.Responder
	onEndState func() //Ending function to close the game
}

type CardPlayer struct {
	ID   string
	hand internal.Deck
	bank int
}

//New generates a new BLKJK
func NewBLKJK() *BLKJK {
	self := &BLKJK{
		name:        "Blackjack",
		description: "Classic casino game of 21",
		rules:       "Standard casino blackjack. The dealer hits on anything under 17.",
		players:     []*CardPlayer{},
		dealer:      internal.NewHand(),
		cursor:      0,
	}

	self.actions = map[string]command.Responder{
		"hit":  self.Hit,
		"stay": self.Stay,
	}

	return self
}

//Name returns the name of this game
func (self *BLKJK) Name() string {
	return self.name
}

//Description returns the description of this game
func (self *BLKJK) Description() string {
	return self.description
}

//Rules returns the rules of this game
func (self *BLKJK) Rules() string {
	return self.rules
}

func (self *BLKJK) Start(s *discordgo.Session, m *discordgo.MessageCreate, finish func()) bool {
	if !self.playerExists(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Please join before starting the game!")
		return false
	}

	command.RegisterDirectCommands(self.actions)
	self.onEndState = finish

	ps := ""
	for _, p := range self.players {
		ps = fmt.Sprintf("%s, <@%s>", ps, p.ID)
	}

	message := fmt.Sprintf("Blackjack has started with the following players: %s", ps[1:])
	s.ChannelMessageSend(m.ChannelID, message)

	self.deck = internal.NewDeck()
	self.deck.Shuffle()
	self.Deal()

	message = "\n[Dealers Hand: " + self.dealerHand(false) + "]"
	for _, p := range self.players {
		message += "\n[<@" + p.ID + ">'s Hand: " + getHand(p.hand) + "]"
	}
	message += "\n<@" + self.players[self.cursor].ID + "> It is your turn."
	s.ChannelMessageSend(m.ChannelID, message)

	return true
}

//Join allows a player to join the session
func (self *BLKJK) Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.players = append(self.players, &CardPlayer{
		ID:   m.Author.ID,
		hand: internal.NewHand(),
		bank: 100,
	})
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has joined Blackjack! There are now %d players", m.Author.ID, len(self.players)))
}

func (self *BLKJK) Deal() {
	for i := 0; i < 2; i++ {
		self.dealer.Insert(self.deck.Draw())
		for _, p := range self.players {
			p.hand.Insert(self.deck.Draw())
		}
	}
}

func (self *BLKJK) Hit(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	current := self.players[self.cursor]
	if m.Author.ID == current.ID {
		current.hand.Insert(self.deck.Draw())
		s.ChannelMessageSend(m.ChannelID, "<@"+current.ID+"> now has "+getHand(current.hand))
		if handTotal(current.hand) > 21 {
			s.ChannelMessageSend(m.ChannelID, "\nOver 21! Bust!")
			self.nextPlayer(s, m)
		}
	}
}

func (self *BLKJK) Stay(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	current := self.players[self.cursor]
	if m.Author.ID == current.ID {
		s.ChannelMessageSend(m.ChannelID, "Staying at "+getHand(current.hand))
		self.nextPlayer(s, m)
	}
}

func (self *BLKJK) nextPlayer(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.cursor++
	message := ""
	if self.cursor == len(self.players) {
		message += "It is now the dealer's turn. The dealer has " + self.dealerHand(true)
		for handTotal(self.dealer) < 17 {
			self.dealer.Insert(self.deck.Draw())
			message += "\nThe dealer draws. The dealer now has " + self.dealerHand(true)
		}

		if handTotal(self.dealer) > 21 {
			message += "\nBust! Congrats to all the winners"
		} else {
			for _, p := range self.players {
				total := handTotal(p.hand)
				if handTotal(self.dealer) >= total || total > 21 {
					message += "\nDealer hand beats <@" + p.ID + ">"
				} else {
					message += "\n<@" + p.ID + "> beats the dealer"
				}
			}
		}
		s.ChannelMessageSend(m.ChannelID, message)
		self.Finish(s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, "\nIt is now <@"+self.players[self.cursor].ID+">'s turn!")
	}
}

func handTotal(deck internal.Deck) int {
	card := deck.Next()
	total := 0
	for card != nil {
		total += cardValue(card, total)
		card = deck.Next()
	}
	return total
}

func (self *BLKJK) dealerHand(full bool) string {
	if !full {
		card := self.dealer.Peek()
		return card.Value() + card.Suit()
	}
	return getHand(self.dealer)
}

func getHand(deck internal.Deck) string {
	message := ""
	c := deck.Next()
	for c != nil {
		message += ", " + c.Value() + c.Suit()
		c = deck.Next()
	}
	return message[2:]
}

func (self *BLKJK) playerExists(ID string) bool {
	for _, p := range self.players {
		if p.ID == ID {
			return true
		}
	}
	return false
}

func cardValue(card internal.Card, total int) int {
	i, e := strconv.Atoi(card.Value())
	if e == nil {
		return i
	}
	switch card.Value() {
	case "J":
		fallthrough
	case "Q":
		fallthrough
	case "K":
		return 10
	case "A":
		if total+11 > 21 {
			return 1
		}
		return 11
	}
	// error
	return 0
}

func (self *BLKJK) Finish(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.players = []*CardPlayer{}
	command.UnregisterDirectCommands(self.actions)
	self.onEndState()
}
