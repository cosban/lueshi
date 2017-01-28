package session

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/gohst/utils"
	ini "github.com/vaughan0/go-ini"
)

var (
	self    *discordgo.User
	session *discordgo.Session
	oauth   *oauth2.Config
	url     string
	secret  string
	debug   bool
)

func init() {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	secret, _ = config.Get("bot", "secret")

	host, _ := config.Get("web", "host")
	_, debug = config.Get("bot", "debug")
	if debug {
		port, _ := config.Get("web", "port")
		url = fmt.Sprintf("http://%s:%s", host, port)
	} else {
		url = fmt.Sprintf("https://%s", host)
	}
}

func Create(user *discordgo.User, discord *discordgo.Session) {
	self = user
	session = discord

	oauth = &oauth2.Config{
		ClientID:     self.ID,
		ClientSecret: secret,
		Scopes:       []string{"identify", "guilds"},
		RedirectURL:  fmt.Sprintf("%s/api/confirm_login", url),
		Endpoint: oauth2.Endpoint{
			AuthURL:  discordgo.EndpointOauth2 + "authorize",
			TokenURL: discordgo.EndpointOauth2 + "token",
		},
	}
}

func Login(w http.ResponseWriter, r *http.Request) error {
	url := oauth.AuthCodeURL("test", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func ConfirmLogin(w http.ResponseWriter, r *http.Request) error {
	code := r.FormValue("code")
	ctx := context.Background()
	token, err := oauth.Exchange(ctx, code)
	if err != nil {
		return err
	}
	u := LoginUser(token)

	session, _ := store.Get(r, "session-name")
	receipt := utils.RandomString(64)
	// perhaps a double map instead...
	for k, v := range sessionIDs {
		if u.ID == v {
			delete(sessionIDs, k)
			break
		}
	}
	sessionIDs[receipt] = u.ID
	session.Values["receipt"] = receipt
	session.Save(r, w)

	UpdateSession(u, receipt)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}
