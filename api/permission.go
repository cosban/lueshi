package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cosban/data"
)

var (
	permissions Set
)

type Set map[int]map[string]bool

func (s Set) Add(discordid string, role int) {
	if _, ok := s[role]; !ok {
		s[role] = make(map[string]bool)
	}
	s[role][discordid] = true
}

func (s Set) Contains(discordid string, role int) bool {
	if b, ok := s[role]; ok {
		p, ok := b[discordid]
		return p && ok
	}
	return false
}

func IsOwner(u *discordgo.User) bool {
	return permissions.Contains(u.ID, 1)
}

func IsTrusted(u *discordgo.User) bool {
	return permissions.Contains(u.ID, 1) || permissions.Contains(u.ID, 2)
}

func IsBanned(u *discordgo.User) bool {
	return permissions.Contains(u.ID, 4)
}

func PopulateRoles() {
	permissions = make(Set)
	rows, err := data.PrepareAndQuery(`SELECT user_id, role_id FROM user_roles`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var discordID string
		var roleID int
		rows.Scan(&discordID, &roleID)
		permissions.Add(discordID, roleID)
	}
}
