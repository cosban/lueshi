package command

import "github.com/bwmarrin/discordgo"

const VERSION = "0.0.3"

func Version(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "version: "+VERSION)
}
