package handler

import "github.com/upfluence/sensu-handler-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/event"

var Store = make(map[string]Handler)

type Handler interface {
	Handle(*event.Event) error
}
