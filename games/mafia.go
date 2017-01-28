package games

import "github.com/bwmarrin/discordgo"

type Mafia struct {
	players                  map[string]bool
	turn                     map[string]bool
	onEndState               func()
	rules, name, description string
}

func NewMafia() *Mafia {
	self := &Mafia{
		players:     map[string]bool{},
		turn:        map[string]bool{},
		name:        "Mafia",
		description: "It's a social deduction game where everyone dies",
		rules:       "Don't get kilt",
	}
	return self
}

func (self *Mafia) Name() string {
	return self.name
}

func (self *Mafia) Description() string {
	return self.description
}

func (self *Mafia) Rules() string {
	return self.rules
}

func (self *Mafia) Start(s *discordgo.Session, m *discordgo.MessageCreate, finish func()) bool {
	return false
}

func (self *Mafia) Finish(*discordgo.Session, *discordgo.MessageCreate) {

}

func (self *Mafia) Join(*discordgo.Session, *discordgo.MessageCreate) {

}
