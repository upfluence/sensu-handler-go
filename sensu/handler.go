package sensu

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/upfluence/sensu-go/sensu/transport"
	"github.com/upfluence/sensu-handler-go/sensu/handler"
)

const (
	connectionTimeout = 5 * time.Second
	currentVersion    = "0.0.1"
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
			log.Printf("Signal %s received", s.String())

			for _, processor := range h.processors {
				processor.Close()
			}

			return h.transport.Close()
		case <-h.transport.GetClosingChan():
			log.Println("transport disconnected")

			for _, processor := range h.processors {
				processor.Close()
			}

			h.transport.Close()
		}
	}
}
