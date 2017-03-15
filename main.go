package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/cosban/lueshi/api"
	"github.com/cosban/lueshi/command"
	"github.com/cosban/lueshi/games"
	"github.com/cosban/lueshi/plugins"
	"github.com/cosban/lueshi/web"

	"github.com/bwmarrin/discordgo"
	"github.com/vaughan0/go-ini"
)

var hearts = []string{
	":sparkling_heart::gift_heart::revolving_hearts::heart_decoration::love_letter:~k-kawaiiiiiiiiii!!!!!~!!!~~!!~:love_letter::heart_decoration::revolving_hearts::gift_heart::sparkling_heart:",
	"<3",
	":heart:",
	"ily %n",
	"go fuck yourself",
	"tee hee, senpai noticed me",
	"わたしは、あなたを愛しています",
	"사랑해",
	"杀了自己",
	":heart:",
	"01101100011011110111011001100101",
	"`" + `func heart(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(hearts)))
	s.ChannelMessageSend(m.ChannelID, strings.ReplaceAll(heart[index], "%n", m.Author.Username, -1))
}` + "`",
}

var (
	token, prefix string
	conf          ini.File
	user          *discordgo.User

	owners = make(map[string]bool)
)

func main() {
	conf = configure()
	discord, err := discordgo.New(token)
	if err != nil {
		log.Fatalln("Failed to create discord session: ", err)
	}
	if user, err = discord.User("@me"); err != nil {
		log.Println("Unable to retrieve credentials: ", err)
	}

	discord.AddHandler(handleMessage)

	if err = discord.Open(); err != nil {
		log.Fatal("Failed to open connection: ", err)
	}
	web.Start(user, discord)
}

func configure() ini.File {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}

	token, _ = config.Get("bot", "token")
	prefix, _ = config.Get("bot", "prefix")
	log.Print("Done. Registering Commands...")
	registerCommands()
	api.PopulateRoles()
	return config
}

func registerCommands() {
	command.RegisterCommand("version", "", "To see the bot's version", command.Version)
	command.RegisterCommand("8ball", "[question]", "Summons the power of a magic eight ball to answer your question (NOTE: often lies)", command.EightBall)
	command.RegisterCommand("yt", "[query]", "Searches youtube for videos", command.Youtube)
	command.RegisterCommand("g", "[query]", "Searches google", command.Google)
	command.RegisterCommand("we", "<city> [state]", "Queries for weather in a given city (SEE ALSO: temp)", command.Weather)
	command.RegisterCommand("temp", "<city> [state]", "Queries for temperature in a given city (SEE ALSO: weather)", command.Temperature)
	command.RegisterCommand("help", "[command]", "Displays these words", command.Help)
	command.RegisterCommand("pasta", "<id | title>", "Craps out some copypasta", plugins.RunPasta)
	command.RegisterCommand("drama", "", "Gives you the latest dramalinks", plugins.NewDrama().Run)
	command.RegisterCommand("quit", "", "Forces the bot to quit", command.Quit)
	command.RegisterCommand("restart", "", "Restarts the bot", command.Restart)
	command.RegisterCommand("set", "<options>", "sets some options for the server", command.Set)
	command.RegisterDirectCommands(games.Commands)
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == user.ID {
		return
	}
	log.Print(m.Content)
	if m.Content == "<3" || m.Content == "❤" {
		heart(s, m)
		return
	}
	content := strings.ToLower(strings.Trim(m.Content, " "))
	args := strings.Split(content, " ")
	channel, _ := s.State.Channel(m.ChannelID)
	if api.IsDirectChannel(channel) {
		if len(args) > 0 {
			command.RunDirectCommand(args[0], args[1:], s, m)
		}
	} else if strings.HasPrefix(content, "<@"+user.ID+">") {
		if len(args) >= 2 {
			command.RunCommand(args[1], args[2:], s, m)
		}
	} else if strings.HasPrefix(content, prefix) {
		if len(args) >= 1 {
			command.RunCommand(args[0][1:], args[1:], s, m)
		}
	}
}

func heart(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(hearts)))
	s.ChannelMessageSend(m.ChannelID, strings.Replace(hearts[index], "%n", m.Author.Username, -1))
}
