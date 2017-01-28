package command

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var responses = []string{
	"It is certain",
	"Don't count on it",
	"It is decidedly so",
	"My reply is no",
	"Without a doubt",
	"My sources say no",
	"Yes, definitely",
	"Outlook not so good",
	"You may rely on it",
	"Very doubtful",
	"As I see it, yes",
	"Definitely not",
	"Mos defs",
	"Hell na",
	"Fuck if I know.",
}

func EightBall(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) < 1 {
		s.ChannelMessageSend(m.ChannelID, "i can't hear you, say it again.")
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(responses)))
	response := fmt.Sprintf("\u200B<@%s>: %s", m.Author.ID, responses[index])
	s.ChannelMessageSend(m.ChannelID, response)
}
