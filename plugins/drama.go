package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/bluemonday"
	ini "github.com/vaughan0/go-ini"
)

type Drama struct {
	username string
	password string
	client   http.Client
	Parse    struct {
		Text struct {
			Content string `json:"*"`
		}
		Externallinks []string
	}
}

func NewDrama() *Drama {
	cookieJar, _ := cookiejar.New(nil)
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	username, _ := config.Get("eti", "username")
	password, _ := config.Get("eti", "password")
	drama := &Drama{
		username: username,
		password: password,
		client: http.Client{
			Jar: cookieJar,
		},
	}
	drama.login()
	return drama
}

func (self *Drama) Run(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	drama := self.requestDrama()
	response := fmt.Sprintf("%s", drama)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "dramalinks",
		Description: response,
		URL:         "https://cosban.net",
	})
}

func (self *Drama) login() {
	data := url.Values{}
	data.Add("b", self.username)
	data.Add("p", self.password)
	r, _ := http.NewRequest("POST", "https://endoftheinter.net", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "*/*")

	resp, err := self.client.Do(r)

	if err != nil {
		log.Println("Error logging into ETI: ", err)
	}
	defer resp.Body.Close()
}

func (self *Drama) requestDrama() string {
	req, err := http.NewRequest("GET", "https://wiki.endoftheinter.net/api.php?action=parse&page=Dramalinks/current&format=json", nil)
	if err != nil {
		panic(err)
	}
	resp, _ := self.client.Do(req)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(content, self)
	return self.sanitize(self.Parse.Text.Content)
}

func (self *Drama) sanitize(drama string) string {
	sanitizer := bluemonday.NewPolicy()
	drama = sanitizer.Sanitize(drama)
	drama = strings.TrimSpace(drama)
	drama = strings.Replace(drama, "&#39;", "'", -1)
	drama = strings.Replace(drama, " update!", "", -1)
	drama = strings.Replace(drama, "\n\n", "", -1)
	drama = drama[:strings.LastIndex(drama, "\n")-2]
	drama = strings.Replace(drama, "Level:", "\nLevel:", 1)
	for i, l := range self.Parse.Externallinks[:len(self.Parse.Externallinks)-1] {
		drama = strings.Replace(drama, fmt.Sprintf("[%d]", i+1), fmt.Sprintf("[[%d]](%s)", i+1, l), -1)
	}
	return drama
}
