package libamabot

import (
	"errors"
	"path/filepath"
	"runtime"
)

// Get the list of all handlers
func GetAllHandlersList() (list []string) {
	for id := range handlers_db {
		list = append(list, id)
	}
	return
}

// Get help of each handler
func GetHelp(id string) (string, error) {
	if h := handlers_db[id]; h != nil {
		return h.help, nil
	}
	return "", errors.New("handler not found named: " + id)
}

// handler to catch request
type handler struct {
	// Handler function. This is to be passed to discordgo.Session through method: AddHandler.
	main interface{}
	// Help message of the handler
	help string
}

// The database that all handlers are in.
// DO NOT EDIT DIRECTRY
var handlers_db map[string]*handler

// add handler to handlers_db when initialize.
func addHandler(h *handler) {
	if handlers_db == nil {
		handlers_db = make(map[string]*handler, 0)
	}
	_, name, _, _ := runtime.Caller(1) // Get the module filename
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])
	handlers_db[name] = h
}
