package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	ini "github.com/vaughan0/go-ini"

	"github.com/cosban/data"
	"github.com/cosban/lueshi/api"
	"github.com/cosban/lueshi/plugins"
	"github.com/cosban/lueshi/web/session"
)

type getlistener func(http.ResponseWriter, *http.Request) error
type postlistener func(http.ResponseWriter, *http.Request) error

type page struct {
	ID                   string
	Title, Content, Path string
	Template             *template.Template
	Info                 interface{}
	User                 *api.User
}

var (
	connection, port string
	self             *discordgo.User
	pageMap          = make(map[string]*page)
)

func init() {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	port, _ = config.Get("web", "port")

	host, _ := config.Get("psql", "host")
	user, _ := config.Get("psql", "username")
	database, _ := config.Get("psql", "database")
	password, _ := config.Get("psql", "password")
	dbport, _ := config.Get("psql", "port")
	log.Print("Attempting to connect to the database...")
	data.Connect(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, dbport, database))
	session.RefreshSessions()
}

func Start(user *discordgo.User, s *discordgo.Session) {
	self = user
	session.Create(user, s)
	listenGet("/", servepage)
	listenAPIGet("/version", version)
	listenAPIGet("/login", session.Login)
	listenAPIGet("/confirm_login", session.ConfirmLogin)
	listenPost("/create_pasta", plugins.HandlePastaCreation)
	log.Printf("Starting web listener on port '%s'", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), nil))
}

func listenAPIGet(uri string, handler getlistener) {
	http.Handle("/api"+uri, handler)
}

func listenGet(uri string, handler getlistener) {
	http.Handle(uri, handler)
}

func listenPost(uri string, handler postlistener) {
	http.Handle("/api"+uri, handler)
}

func (fn postlistener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "METHOD NOT ALLOWED", 405)
		return
	}
	if !session.IsConnected(r) {
		http.Error(w, "NOT AUTHORIZED", 401)
		return
	}
	if err := fn(w, r); err != nil {
		log.Print(err)
		http.Error(w, "ERROR", 500)
	}
}

func (fn getlistener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Print(err)
		http.Error(w, "ERROR", 500)
	}
}

func servepage(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "METHOD NOT ALLOWED", 405)
		return nil
	}
	title := r.URL.Path[len("/"):]
	if strings.HasSuffix(title, "/") {
		title = title[:len(title)-1]
	}
	if len(title) < 1 {
		title = "home"
	}
	p := loadpage(w, title)
	if p == nil {
		http.Error(w, "NOT FOUND", 404)
		return nil
	}
	if session.IsConnected(r) {
		p.User = session.GetConnectedUser(r)
	}
	if f, ok := info[title]; ok {
		p.Info = f(w, r)
	}
	return p.Template.ExecuteTemplate(w, "base.html", p)
}

func loadpage(w http.ResponseWriter, title string) *page {
	filename := "pages/" + title + ".html"
	if strings.Contains(title, "/") {
		title = title[strings.LastIndex(title, "/")+1:]
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	t := template.Must(template.ParseFiles("pages/base.html", filename))
	return &page{
		Title:    title,
		Template: t,
		Content:  string(body),
		ID:       self.ID,
	}
}
