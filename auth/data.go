package auth

import (
	"github.com/AskUbuntu/tbot/util"

	"sync"
)

type data struct {
	sync.Mutex
	name  string
	Users map[string]*User `json:"users"`
}

func (d *data) load() error {
	_, err := util.LoadJSON(d.name, d)
	return err
}

func (d *data) save() error {
	return util.SaveJSON(d.name, d)
}
