package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/api"
)

func Set(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !api.IsOwner(m.Author) {
		return
	}

}
