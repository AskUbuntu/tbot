package server

import (
	"github.com/AskUbuntu/tbot/queue"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"
	"github.com/gorilla/Schema"

	"net/http"
	"strings"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "index.html", pongo2.Context{})
}

func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		scraperSettings = s.scraper.Settings()
		queueSettings   = s.queue.Settings()
		f               = struct {
			ScraperSettings scraper.Settings
			QueueSettings   queue.Settings
			MatchingWords   string
		}{
			scraperSettings,
			queueSettings,
			strings.Join(scraperSettings.MatchingWords, ", "),
		}
	)
	if r.Method == http.MethodPost {
		decoder := schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)
		if err := decoder.Decode(&f, r.Form); err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			f.ScraperSettings.MatchingWords = util.SplitAndTrimString(
				f.MatchingWords, ",",
			)
			if err := util.AnyError(
				s.scraper.SetSettings(f.ScraperSettings),
				s.queue.SetSettings(f.QueueSettings),
			); err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				s.addAlert(w, r, infoType, "settings saved to disk")
				http.Redirect(w, r, "settings", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "settings.html", pongo2.Context{
		"title": "Settings",
		"f":     f,
	})
}
