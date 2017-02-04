package games

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"strconv"

	"github.com/cosban/lueshi/command"
	"github.com/cosban/lueshi/internal"
)

type CardPlayer struct {
	ID   string
	hand internal.Deck
	bank int
}

func NewPlayer(id string) *CardPlayer {
	return &CardPlayer{
		ID:   id,
		hand: internal.NewHand(),
		bank: 100,
	}
}

type BLKJK struct {
	name        string
	description string
	rules       string
	players     []*CardPlayer
	bets        map[string]int
	dealer      internal.Deck
	deck        internal.Deck
	cursor      int
	betting     bool

	actions    map[string]command.Responder
	onEndState func()
}

func NewBLKJK() *BLKJK {
	self := &BLKJK{
		name:        "Blackjack",
		description: "Classic casino game of 21",
		rules:       "Standard casino blackjack. The dealer hits on anything under 17.",
		players:     []*CardPlayer{},
		bets:        map[string]int{},
		dealer:      internal.NewHand(),
		cursor:      0,
		betting:     true,
	}

	self.actions = map[string]command.Responder{
		"hit":     self.Hit,
		"stay":    self.Stay,
		"bet":     self.Bet,
		"leave":   self.Leave,
		"balance": self.Total,
	}

	return self
}

func (self *BLKJK) Name() string {
	return self.name
}

func (self *BLKJK) Description() string {
	return self.description
}

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
	message += "\nPlease make your bets now."
	s.ChannelMessageSend(m.ChannelID, message)

	self.betting = true
	return true
}

func (self *BLKJK) Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := NewPlayer(m.Author.ID)
	if self.addPlayer(player) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has joined Blackjack! There are now %d players", m.Author.ID, len(self.players)))
	} else {
		s.ChannelMessageSend(m.ChannelID, "You are already ready to play!")
	}
}

func (self *BLKJK) Leave(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	player := self.getPlayer(m.Author.ID)
	if self.removePlayer(player) {
		s.ChannelMessageSend(m.ChannelID, "come back soon.")
	}
	if len(self.players) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Also no one is left so the game is finished.")
		self.Finish(s, m)
	}
}

func (self *BLKJK) Finish(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.players = []*CardPlayer{}
	self.dealer = internal.NewHand()
	self.cursor = 0
	self.deck = internal.NewDeck()

	command.UnregisterDirectCommands(self.actions)
	self.onEndState()
}

func (self *BLKJK) addPlayer(player *CardPlayer) bool {
	if self.playerExists(player.ID) {
		return false
	}
	self.players = append(self.players, player)
	return true
}

func (self *BLKJK) removePlayer(player *CardPlayer) bool {
	for i, p := range self.players {
		if p.ID == player.ID {
			self.players = append(self.players[:i], self.players[i+1:]...)
			delete(self.bets, p.ID)
			return true
		}
	}
	return false
}

func (self *BLKJK) dealNewHand() {
	self.deck = internal.NewDeck()
	self.deck.Shuffle()
	for i := 0; i < 2; i++ {
		self.dealer.Insert(self.deck.Draw())
		for _, p := range self.players {
			p.hand.Insert(self.deck.Draw())
		}
	}
}

func (self *BLKJK) dealHand(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.dealNewHand()
	message := fmt.Sprintf("\n[Dealer: %s]", self.dealerHand(false))
	for _, p := range self.players {
		message += fmt.Sprintf("\n[<@%s>: %s] (%d)", p.ID, getHand(p.hand), handTotal(p.hand))
	}
	message += "\n<@" + self.players[self.cursor].ID + "> It is your turn."
	s.ChannelMessageSend(m.ChannelID, message)
}

func (self *BLKJK) hit(id string) bool {
	current := self.players[self.cursor]
	if id == current.ID {
		current.hand.Insert(self.deck.Draw())
		return true
	}
	return false
}

func (self *BLKJK) stay(id string) bool {
	current := self.players[self.cursor]
	if id == current.ID {
		self.cursor++
		return true
	}
	return false
}

func (self *BLKJK) Bet(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !self.betting {
		s.ChannelMessageSend(m.ChannelID, "Betting is not allowed at this time")
	}
	p := self.getPlayer(m.Author.ID)
	if p == nil {
		return
	}
	message := ""
	if b, e := strconv.Atoi(args[0]); e == nil && b <= p.bank && b > 1 {
		self.bets[p.ID] = b
		p.bank = p.bank - b
		message += fmt.Sprintf("<@%s> bets %d and has %d remaining", p.ID, b, p.bank)
	} else {
		message += fmt.Sprintf("<@%s> you can't bet that. You have %d remaining", p.ID, p.bank)
		s.ChannelMessageSend(m.ChannelID, message)
		return
	}
	if len(self.bets) == len(self.players) {
		self.betting = false
		s.ChannelMessageSend(m.ChannelID, message+"\nBetting is now closed")
		self.dealHand(s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, message)
	}
}

func (self *BLKJK) Total(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, p := range self.players {
		if p.ID == m.Author.ID {
			message := fmt.Sprintf("<@%s> has %d yetibux.", p.ID, p.bank)
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}

func (self *BLKJK) Hit(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if self.hit(m.Author.ID) {
		p := self.players[self.cursor]
		message := fmt.Sprintf("\n[<@%s>: %s] (%d)", p.ID, getHand(p.hand), handTotal(p.hand))
		if handTotal(p.hand) > 21 {
			s.ChannelMessageSend(m.ChannelID, message+"\nOver 21! Bust!")
			self.nextPlayer(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}

func (self *BLKJK) Stay(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	p := self.players[self.cursor]
	if self.stay(m.Author.ID) {
		message := fmt.Sprintf("\nStay [<@%s>: %s] (%d)", p.ID, getHand(p.hand), handTotal(p.hand))
		s.ChannelMessageSend(m.ChannelID, message)
		self.nextPlayer(s, m)
	}
}

func (self *BLKJK) nextPlayer(s *discordgo.Session, m *discordgo.MessageCreate) {
	if self.cursor < len(self.players) {
		s.ChannelMessageSend(m.ChannelID, "\nIt is now <@"+self.players[self.cursor].ID+">'s turn!")
		return
	}
	message := "It is now the dealer's turn. The dealer has " + self.dealerHand(true)
	for handTotal(self.dealer) < 17 {
		self.dealer.Insert(self.deck.Draw())
		message += "\nThe dealer draws and now has " + self.dealerHand(true)
	}
	if handTotal(self.dealer) > 21 {
		message += "Dealer bust!"
	}
	message += self.dealEarnings(s, m)
	if len(self.players) > 0 {
		s.ChannelMessageSend(m.ChannelID, message+"\nA new round is now starting. Please make your bets.")
		self.reset()
	} else {
		s.ChannelMessageSend(m.ChannelID, message+"\nLol you guys fucking suck. The game is now over.")
		self.Finish(s, m)
	}
}

func (self *BLKJK) dealEarnings(s *discordgo.Session, m *discordgo.MessageCreate) string {
	message := ""
	for _, p := range self.players {
		total := handTotal(p.hand)
		dealer := handTotal(self.dealer)
		if total <= 21 && total > dealer || total <= 21 && dealer > 21 {
			p.bank = p.bank + self.bets[p.ID]*2
			message += fmt.Sprintf("\n<@%s> beats the dealer and earns %d (bank: %d)", p.ID, self.bets[p.ID]*2, p.bank)
		} else {
			message += fmt.Sprintf("\nDealer hand beats <@%s> (bank: %d)", p.ID, p.bank)
			if p.bank == 0 {
				message += "... and because they have no more money I'm kicking them out."
				self.removePlayer(p)
			}
		}
	}
	return message
}

func (self *BLKJK) reset() {
	for _, p := range self.players {
		p.hand = internal.NewHand()
	}
	self.dealer = internal.NewHand()
	self.bets = map[string]int{}
	self.cursor = 0
	self.deck = internal.NewDeck()
	self.betting = true
}

func (self *BLKJK) dealerHand(full bool) string {
	if !full {
		card := self.dealer.Peek()
		return card.Value() + card.Suit()
	}
	return getHand(self.dealer)
}

func (self *BLKJK) getPlayer(ID string) *CardPlayer {
	for _, p := range self.players {
		if p.ID == ID {
			return p
		}
	}
	return nil
}

func (self *BLKJK) playerExists(ID string) bool {
	return self.getPlayer(ID) != nil
}

func getHand(deck internal.Deck) string {
	message := ""
	if deck.Size() < 1 {
		return message
	}
	c := deck.Next()
	for c != nil {
		message += ", " + c.Value() + c.Suit()
		c = deck.Next()
	}
	return message[2:]
}

func handTotal(deck internal.Deck) int {
	card := deck.Next()
	total := 0
	a := 0
	for card != nil {
		total, a = addCardValue(card, a, total)
		card = deck.Next()
	}
	return total
}

// card: a card
// a: total of maxed aces
// total: the current hand total
func addCardValue(card internal.Card, a, total int) (int, int) {
	i, e := strconv.Atoi(card.Value())
	if e != nil {
		switch card.Value() {
		case "J":
			fallthrough
		case "Q":
			fallthrough
		case "K":
			i = 10
		case "A":
			i = 11
			a++
		}
	}
	total = total + i
	for a > 0 && total > 21 {
		a--
		total = total - 10
	}
	return total, a
}
