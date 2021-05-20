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
	if h := help_db[id]; h != "" {
		return h, nil
	}
	return "", errors.New("help not found named: " + id)
}

// The database that all handlers are in.
// DO NOT EDIT DIRECTRY
var handlers_db map[string][]interface{}
var help_db map[string]string

// add handler to handlers_db when initialize.
func addHandler(h ...interface{}) {
	if handlers_db == nil {
		handlers_db = make(map[string][]interface{}, 0)
	}
	_, name, _, _ := runtime.Caller(1) // Get the module filename
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])
	handlers_db[name] = h
}

// add help to help_db when initialize.
func addHelp(h ...interface{}) {
	if handlers_db == nil {
		handlers_db = make(map[string][]interface{}, 0)
	}
	_, name, _, _ := runtime.Caller(1) // Get the module filename
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])
	handlers_db[name] = h
}
