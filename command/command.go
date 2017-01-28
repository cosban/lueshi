package command

import (
	"log"

	"github.com/bwmarrin/discordgo"
	ini "github.com/vaughan0/go-ini"
)

var (
	key, cx, weatherkey string
	commands            = make(map[string]Command)
	direct              = make(map[string]Command)
)

type Responder func([]string, *discordgo.Session, *discordgo.MessageCreate)

type Command struct {
	Name, Usage, Summary string
	Respond              Responder
}

func init() {
	conf, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	key, _ = conf.Get("google", "key")
	cx, _ = conf.Get("google", "cx")
	weatherkey, _ = conf.Get("weather", "key")
}

func RegisterCommand(name, usage, summary string, responder Responder) {
	commands[name] = Command{name, usage, summary, responder}
}

func RunCommand(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if c, ok := commands[command]; ok {
		c.Respond(args, s, m)
	}
}

func RegisterDirectCommands(commands map[string]Responder) {
	for k, v := range commands {
		direct[k] = Command{k, "", "", v}
	}
}

func UnregisterDirectCommands(commands map[string]Responder) {
	for k := range commands {
		delete(direct, k)
	}
}

func RunDirectCommand(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if c, ok := direct[command]; ok {
		c.Respond(args, s, m)
	}
}
