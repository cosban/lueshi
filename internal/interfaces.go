package internal

import "github.com/bwmarrin/discordgo"

type Plugin interface {
	Run(args []string, s *discordgo.Session, m *discordgo.MessageCreate)
}

// Game is a simple interface which allows you to add "game modules"
type Game interface {
	// Name returns the name of the module
	Name() string
	// Description returns the description of the module
	Description() string
	// Rules returns the specific rules of the module
	Rules() string
	// Start takes an "onFinish" function as well as the normal session and event params
	// It is called by the handler when the game starts
	Start(*discordgo.Session, *discordgo.MessageCreate, func()) bool
	// Finish should perform cleanup tasks and then call the "onFinish" function which was passed in at Start
	Finish(*discordgo.Session, *discordgo.MessageCreate)
	// Join is called when a player joins the game
	Join(*discordgo.Session, *discordgo.MessageCreate)
}

type Deck interface {
	Shuffle()
	Insert(Card)
	Draw() Card
	Peek() Card
	Next() Card
	Size() int
}

type Card interface {
	Value() string
	Suit() string
}
