package sensu

import (
	"os"
	"os/signal"
	"time"

	"github.com/upfluence/sensu-handler-go/Godeps/_workspace/src/github.com/upfluence/goutils/log"
	"github.com/upfluence/sensu-handler-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/transport"
	"github.com/upfluence/sensu-handler-go/sensu/handler"
)

const (
	connectionTimeout = 5 * time.Second
	currentVersion    = "0.1.0"
)

type Handler struct {
	transport  transport.Transport
	processors []*processor
}

func NewHandler(transport transport.Transport) *Handler {
	processors := []*processor{}

	for name, handler := range handler.Store {
		processors = append(processors, &processor{name, handler, transport, nil})
	}

	return &Handler{transport, processors}
}

func (h *Handler) Start() error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	for {
		h.transport.Connect()

		for !h.transport.IsConnected() {
			select {
			case <-time.After(connectionTimeout):
				h.transport.Connect()
			case <-sig:
				return h.transport.Close()
			}
		}

		for _, processor := range h.processors {
			go processor.Start()
		}

		select {
		case s := <-sig:
			log.Noticef("Signal %s received", s.String())

			for _, processor := range h.processors {
				processor.Close()
			}

			return h.transport.Close()
		case <-h.transport.GetClosingChan():
			log.Warning("transport disconnected")

			for _, processor := range h.processors {
				processor.Close()
			}

			h.transport.Close()
		}
	}
}
