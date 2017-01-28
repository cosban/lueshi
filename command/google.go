package command

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/internal"
)

type SearchData struct {
	Items []struct {
		Kind, Title, HTMLTitle, Link, DisplayLink, Snippet, HTMLSnippet, CacheID, FormattedURL, HTMLFormattedURL string
	}
	Error struct {
		Message string
	}
}

// Search performs a google search
func Google(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) < 1 {
		return
	}
	search := strings.Join(args, " ")
	site := parseSite(search)
	q := url.QueryEscape(search)

	request := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?&key=%s&cx=%s&q=%s&siteSearch=%s&fields=items&num=1", key, cx, q, site)
	r := &SearchData{}
	internal.GetJSON(request, r)

	response := fmt.Sprintf("\u200B<@%s>: No results found :(", m.Author.ID)
	if len(r.Items) > 0 {
		response = fmt.Sprintf("\u200B<@%s>: %s -- \u0002%s\u0002: \"%s\" ", m.Author.ID, r.Items[0].Link, r.Items[0].Title, r.Items[0].Snippet)
	} else if len(r.Error.Message) > 0 {
		log.Print(r.Error.Message)
	}
	s.ChannelMessageSend(m.ChannelID, response)
}

func parseSite(s string) string {
	for _, element := range strings.Split(s, " ") {
		if strings.HasPrefix(element, "site:") {
			return element[len("site:"):]
		}
	}
	return ""
}
