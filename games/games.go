package games

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/command"
	"github.com/cosban/lueshi/internal"
)

var (
	Commands = map[string]command.Responder{
		"mode":   mode,
		"start":  start,
		"join":   join,
		"games":  list,
		"status": status,
	}
	running = false
	games   = map[string]internal.Game{
		"attack/dodge": NewATKD(),
		"blackjack":    NewBLKJK(),
		//"mafia":        NewMafia(),
	}
	current internal.Game
)

func mode(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if running {
		s.ChannelMessageSend(m.ChannelID, "Please wait for the current game to finish.")
	} else if len(args) > 0 {
		game := strings.ToLower(args[0])
		selection, ok := games[game]
		if !ok {
			s.ChannelMessageSend(m.ChannelID, "I've never heard of that game.")
		} else {
			current = selection
			message := fmt.Sprintf("%s has been selected", current.Name())
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}

func start(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if running {
		s.ChannelMessageSend(m.ChannelID, "Please wait for the current game to finish")
	} else if current != nil {
		running = current.Start(s, m, finish)
	}
}

func join(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if running {
		s.ChannelMessageSend(m.ChannelID, "Please wait for the current game to finish")
	} else {
		current.Join(s, m)
	}
}

func list(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	response := ""
	for k, v := range games {
		response = fmt.Sprintf("%s\n%s - %s", response, k, v.Description())
	}
	s.ChannelMessageSend(m.ChannelID, response[0:])
}

func rules(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if current != nil {
		rules := current.Rules()
		s.ChannelMessageSend(m.ChannelID, rules)
	} else {
		s.ChannelMessageSend(m.ChannelID, "No game has been selected")
	}
}

func finish() {
	running = false
}

func status(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if current != nil {
		message := fmt.Sprintf("The current game is %s and it is ", current.Name())
		if !running {
			message += "not "
		}
		message += "running"
		s.ChannelMessageSend(m.ChannelID, message)
	} else {
		s.ChannelMessageSend(m.ChannelID, "no game is currently selected")
	}
}
