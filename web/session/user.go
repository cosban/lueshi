package session

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/data"
	"github.com/cosban/lueshi/api"
)

func IsConnected(r *http.Request) bool {
	session, _ := store.Get(r, "session-name")
	if receipt, ok := session.Values["receipt"].(string); ok {
		if _, existsID := sessionIDs[receipt]; existsID {
			return true
		}
	}
	return false
}

func GetConnectedUser(r *http.Request) *api.User {
	s, _ := store.Get(r, "session-name")
	if receipt, ok := s.Values["receipt"].(string); ok {
		if v, _ := sessionIDs[receipt]; len(v) > 0 {
			return api.GetUserFromID(sessionIDs[receipt])
		}
	}
	return nil
}

func LoginUser(token *oauth2.Token) *api.User {
	s, _ := discordgo.New(token.Type() + " " + token.AccessToken)
	u, _ := s.User("@me")
	data.PrepareAndExecute(
		`INSERT INTO users (user_id, username) 
		VALUES ($1, $2) 
		ON CONFLICT(user_id)
		DO UPDATE SET username = $2`, u.ID, u.Username)
	return &api.User{
		ID:       u.ID,
		Username: u.Username,
	}
}
