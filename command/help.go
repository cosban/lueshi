package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Help(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(args) < 1 {
		str := ""

		for _, v := range commands {
			str = fmt.Sprintf("%s\n%s -- %s", str, v.Name, v.Summary)
		}
		s.ChannelMessageSend(m.ChannelID, "More information is available at https://lueshi.cosban.net")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("NOTE: To view usage of a command, please use .help <command>%s", str))
	} else if v, ok := commands[args[0]]; ok {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s %s -- %s", v.Name, v.Usage, v.Summary))
	}
}
