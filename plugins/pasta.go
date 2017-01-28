package plugins

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/data"
	"github.com/cosban/lueshi/web/session"
)

type Pasta struct {
	ID, Server, Title, Content, Creator string
}

type PastaResponse struct {
	Current *Pasta
	List    []*Pasta
	Server  string
}

func RunPasta(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := s.State.Channel(m.ChannelID)
	server := channel.GuildID
	if len(args) < 1 {
		response := getRandomPasta(server)
		s.ChannelMessageSend(m.ChannelID, response.Content)
		return
	}
	response := RetreiveByTitle(strings.Join(args, " "), server)
	s.ChannelMessageSend(m.ChannelID, response.Content)
}

func RetreiveByID(id, server string) *Pasta {
	var title, content, creator string
	err := data.QueryRow(data.Prepare(`
    SELECT title, content, user_id
    FROM pasta
    WHERE id = $1 AND server_id = $2;`, id, server), &title, &content, &creator)
	if err != nil {
		panic(err)
	}
	return &Pasta{
		ID:      id,
		Title:   title,
		Content: content,
		Creator: creator,
	}
}

func RetreiveAll(server string) []*Pasta {
	rows, err := data.PrepareAndQuery(
		`SELECT id, title, content, user_id 
		FROM pasta 
		WHERE server_id = $1`, server)
	if err != nil {
		panic(err)
	}
	var pastas []*Pasta
	for rows.Next() {
		var id, title, content, creator string
		err = rows.Scan(&id, &title, &content, &creator)
		if err != nil {
			panic(err)
		}
		pastas = append(pastas, &Pasta{
			ID:      id,
			Title:   title,
			Content: content,
			Creator: creator,
			Server:  server,
		})
	}
	return pastas
}

func RetreiveByTitle(query, server string) *Pasta {
	var id, title, content, creator string
	err := data.QueryRow(data.Prepare(
		`SELECT id, title, content, user_id 
		FROM pasta
		WHERE LOWER(title) LIKE $1
		AND server_id = $2;`, strings.ToLower(query), server), &id, &title, &content, &creator)
	if err != nil {
		panic(err)
	}
	return &Pasta{
		ID:      id,
		Title:   title,
		Content: content,
		Creator: creator,
		Server:  server,
	}
}

func getRandomPasta(server string) *Pasta {
	pastas := RetreiveAll(server)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Int31n(int32(len(pastas)))
	return &Pasta{
		ID:      pastas[index].ID,
		Title:   pastas[index].Title,
		Content: pastas[index].Content,
		Creator: pastas[index].Creator,
		Server:  server,
	}
}

func createPasta(creator, title, content, server string) (int, error) {
	var id int
	err := data.QueryRow(data.Prepare(
		`INSERT INTO pasta (user_id, title, server_id, content) 
		VALUES($1, $2, $3, $4) 
		RETURNING id;`, creator, title, server, content), &id)
	return id, err
}

func HandlePastaCreation(w http.ResponseWriter, r *http.Request) error {
	u := session.GetConnectedUser(r)
	title := r.FormValue("title")
	content := r.FormValue("content")
	server := r.FormValue("server")
	log.Println(u.Username, "created pasta", title)
	id, err := createPasta(u.ID, title, content, server)
	http.Redirect(w, r, fmt.Sprintf("/pasta?s=%s&id=%d", server, id), http.StatusSeeOther)
	return err
}

func GetPasta(w http.ResponseWriter, r *http.Request) interface{} {
	id := r.FormValue("id")
	server := r.FormValue("s")
	response := PastaResponse{
		List:   RetreiveAll(server),
		Server: server,
	}
	if len(id) > 0 {
		response.Current = RetreiveByID(id, server)
	}
	return response
}
