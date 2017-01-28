package command

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/internal"
)

type YoutubeData struct {
	Items []struct {
		Id struct {
			VideoID string
		}
		Snippet struct {
			Title, Description, ChannelTitle string
		}
	}
}

func Youtube(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(args) < 1 {
		return
	}
	search := strings.Join(args, " ")
	q := url.QueryEscape(search)

	request := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?&key=%s&part=id,snippet&maxResults=1&q=%s", key, q)
	r := &YoutubeData{}
	internal.GetJSON(request, r)

	response := fmt.Sprintf("\u200B<@%s>: No results found :(", m.Author.ID)
	if len(r.Items) > 0 {
		response = fmt.Sprintf(
			"\u200B<@%s>: https://youtube.com/watch?v=%s -- \u0002%s by %s\u0002: \"%s\" ",
			m.Author.ID,
			r.Items[0].Id.VideoID,
			r.Items[0].Snippet.Title,
			r.Items[0].Snippet.ChannelTitle,
			r.Items[0].Snippet.Description,
		)
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
