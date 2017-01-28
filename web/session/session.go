package session

import (
	"github.com/cosban/data"
	"github.com/cosban/lueshi/api"
	"github.com/gorilla/sessions"
	ini "github.com/vaughan0/go-ini"
)

var (
	store      *sessions.CookieStore
	sessionIDs map[string]string
)

func init() {
	config, err := ini.LoadFile("config.ini")
	if config != nil && err == nil {
		key, _ := config.Get("cookies", "store")
		store = sessions.NewCookieStore([]byte(key))
	}
}

func RefreshSessions() {
	if len(sessionIDs) < 1 {
		sessionIDs = GetSessions()
	}
}

func GetSessions() map[string]string {
	session := make(map[string]string)
	stmt := data.Prepare("SELECT user_id, session FROM users WHERE session IS NOT NULL;")
	rows, err := data.Query(stmt)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var userid, receipt string
		rows.Scan(&userid, &receipt)
		session[receipt] = userid
	}
	return session
}

func UpdateSession(u *api.User, receipt string) error {
	return data.PrepareAndExecute(
		`UPDATE users 
        SET session = $1
        WHERE user_id = $2;`, receipt, u.ID)
}

func RemoveSession(userid int) error {
	return data.PrepareAndExecute(
		`UPDATE users 
         SET session = NULL
         WHERE user_id = $1`, userid)
}
