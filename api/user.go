package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cosban/data"
	"golang.org/x/oauth2"
)

type User struct {
	ID       string
	Username string
}

func GetUserFromID(userid string) *User {
	stmt := data.Prepare(`SELECT username, COUNT(username) FROM users WHERE user_id = $1 GROUP BY username;`, userid)
	var username string
	var count int
	data.QueryRow(stmt, &username, &count)
	if count != 1 {
		return nil
	}
	return &User{
		ID:       userid,
		Username: username,
	}
}

func GetUserFromToken(token *oauth2.Token) *User {
	s, _ := discordgo.New(token.Type() + " " + token.AccessToken)
	u, _ := s.User("@me")
	stmt := data.Prepare(
		`SELECT user_id, username
		FROM users
		WHERE discord_id = $1`, u.ID)
	var id, username string
	data.QueryRow(stmt, &id, &username)

	return &User{
		ID:       id,
		Username: username,
	}
}

func GetUserServers(userid string) []*Server {
	rows, _ := data.PrepareAndQuery(
		`SELECT s.name, s.server_id
		FROM server s
		JOIN server_users su USING (server_id)
		WHERE su.user_id = $1;`, userid)
	var servers []*Server
	for rows.Next() {
		var name, id string
		rows.Scan(&name, &id)
		servers = append(
			servers,
			&Server{
				ID:   id,
				Name: name,
			},
		)
	}
	return servers
}
