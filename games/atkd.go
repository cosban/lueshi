package games

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/command"
)

type ATKD struct {
	players                  map[string]bool
	turn                     map[string]bool
	actions                  map[string]command.Responder
	onEndState               func()
	rules, name, description string
}

func NewATKD() *ATKD {
	self := &ATKD{
		players:     map[string]bool{},
		turn:        map[string]bool{},
		name:        "Attack/Dodge",
		description: "It's like rock, paper, scissors but everyone dies",
		rules: "Attack/Dodge is the easiest game requiring two or more people on the planet.\n" +
			"The goal is simple: Kill your opponent by attacking them!\n" +
			"Do this by strategically choosing whether to attack or dodge.\n" +
			"Each move has the following effects:\n\n" +
			"ATTACK -- Allows you to kill anyone that hasn't dodged (except you)\n" +
			"DODGE  -- Prevent your death! If you dodge, you can not be attacked during this round!",
	}
	self.actions = map[string]command.Responder{
		"attack": self.attack,
		"dodge":  self.dodge,
	}

	return self
}

func (self *ATKD) Name() string {
	return self.name
}

func (self *ATKD) Description() string {
	return self.description
}

func (self *ATKD) Rules() string {
	return self.rules
}

func (self *ATKD) Start(s *discordgo.Session, m *discordgo.MessageCreate, finish func()) bool {
	if _, ok := self.players[m.Author.ID]; !ok {
		s.ChannelMessageSend(m.ChannelID, "Please join before starting the game!")
		return false
	}
	command.RegisterDirectCommands(self.actions)
	self.onEndState = finish
	ps := ""
	for p := range self.players {
		ps = fmt.Sprintf("%s, <@%s>", ps, p)
	}
	message := fmt.Sprintf("Attack/Dodge has started with the following players: %s", ps[1:])
	s.ChannelMessageSend(m.ChannelID, message)
	self.turn = map[string]bool{}
	return true
}

func (self *ATKD) Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.players[m.Author.ID] = true
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has joined Attack/Dodge! There are now %d players.", m.Author.ID, len(self.players)))
}

func (self *ATKD) dodge(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	move(self, false, s, m)
}

func (self *ATKD) attack(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	move(self, true, s, m)
}

func move(self *ATKD, attacked bool, s *discordgo.Session, m *discordgo.MessageCreate) {
	if attacked {
		s.ChannelMessageSend(m.ChannelID, "You have chosen to attack!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "You have chosen to dodge!")
	}
	self.turn[m.Author.ID] = attacked
	if len(self.players) == len(self.turn) {
		endTurn(self, s, m)
	}
}

func endTurn(self *ATKD, s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""
	dead := make(map[string]bool)
	for player, attacked := range self.turn {
		if attacked {
			for victim, didAttack := range self.turn {
				if victim != player && didAttack {
					if _, ok := dead[victim]; !ok {
						message += "\n<@" + player + "> has killed " + "<@" + victim + ">"
						dead[victim] = true
					}
				}
			}
		}
	}

	self.turn = map[string]bool{}
	for d := range dead {
		delete(self.players, d)
	}
	if len(self.players) == 1 {
		for p := range self.players {
			message += "\nWE HAVE A WINNER! Congratulations to <@" + p + "> for not dying!"
		}
		s.ChannelMessageSend(m.ChannelID, message)
		self.Finish(s, m)
	} else if len(self.players) == 0 {
		message += "\nWow you guys fucking suck, everyone died!"
		s.ChannelMessageSend(m.ChannelID, message)
		self.Finish(s, m)
	} else {
		if len(dead) == 0 {
			message += "\nNo one died!"
		}
		message += "\nIt is now the start of a new turn!"
		s.ChannelMessageSend(m.ChannelID, message)
	}
}

func (self *ATKD) Finish(s *discordgo.Session, m *discordgo.MessageCreate) {
	self.players = map[string]bool{}
	command.UnregisterDirectCommands(self.actions)
	self.onEndState()
}
