package web

import (
	"fmt"
	"net/http"

	"github.com/cosban/lueshi/api"
	"github.com/cosban/lueshi/command"
	"github.com/cosban/lueshi/plugins"
	"github.com/cosban/lueshi/web/session"
)

var info = map[string]func(http.ResponseWriter, *http.Request) interface{}{
	"pasta":     plugins.GetPasta,
	"dashboard": getDashboard,
}

type dashboardResponse struct {
	Servers []*api.Server
}

func getDashboard(w http.ResponseWriter, r *http.Request) interface{} {
	u := session.GetConnectedUser(r)
	servers := api.GetUserServers(u.ID)
	return dashboardResponse{
		Servers: servers,
	}
}

func version(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		http.Error(w, "METHOD NOT ALLOWED", 405)
		return nil
	}
	fmt.Fprintf(w, command.VERSION)
	return nil
}
