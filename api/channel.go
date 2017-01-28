package api

import (
	"github.com/bwmarrin/discordgo"

	"github.com/cosban/data"
)

// TODO: caching
func IsDirectChannel(channel *discordgo.Channel) bool {
	var count int
	data.QueryRow(data.Prepare(`SELECT COUNT(*) FROM server WHERE direct_channel = $1;`, channel.ID), &count)
	return count > 0
}
