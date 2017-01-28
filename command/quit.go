package command

import (
	"log"
	"math/rand"
	"os"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/api"
)

var quits = []string{
	"*ollies out*",
	"piece bitches",
	"fuck this shit",
	"*kills self*",
	"lol this is gay",
	"*goes back to the topic list*",
	"yikes im having a meltie",
	"this conversation is triggering me",
	"see you space pudding pop",
	"as you wish",
	"sayonara my rolling star",
	"one more once?",
	"saving game... please do not turn of the console or remove the memory card",
	"RIP bankroll fresh",
}

func Quit(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !api.IsOwner(m.Author) {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(quits)))
	s.ChannelMessageSend(m.ChannelID, quits[index])
	os.Exit(1)
}

func Restart(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !api.IsOwner(m.Author) {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(quits)))
	s.ChannelMessageSend(m.ChannelID, quits[index])

	gopath := os.Getenv("GOPATH")
	err := syscall.Exec(gopath+"/bin/lueshi", []string{"lueshi"}, os.Environ())
	log.Print(err)
}
