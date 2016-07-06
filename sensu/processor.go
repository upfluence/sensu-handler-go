package sensu

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/upfluence/sensu-go/sensu/event"
	"github.com/upfluence/sensu-go/sensu/transport"
	"github.com/upfluence/sensu-handler-go/sensu/handler"
)

type processor struct {
	name      string
	handler   handler.Handler
	transport transport.Transport
	closeChan chan bool
}

func (p *processor) Start() error {
	funnel := strings.Join(
		[]string{
			p.name,
			currentVersion,
			strconv.Itoa(int(time.Now().Unix())),
		},
		"-",
	)

	msgChan, stopChan := p.subscribe(funnel)
	log.Printf("Subscribed to %s", p.name)

	for {
		select {
		case b := <-msgChan:
			p.handleMessage(b)
		case <-p.closeChan:
			log.Printf("Gracefull stop of %s", p.name)
			stopChan <- true
			return nil
		}
	}

	return nil
}

func (p *processor) handleMessage(blob []byte) {
	event, err := event.UnmarshalEvent(blob)

	if err != nil {
		log.Printf("%s: unmarshal the event: %s", string(blob), err.Error())
	} else {
		log.Printf("Event received: %s", string(blob))
	}

	if err := p.handler.Handle(event); err != nil {
		log.Printf("%s: %s", p.name, err.Error())
	}
}

func (p *processor) Close() {
	if p.closeChan != nil {
		p.closeChan <- true
	}
}

func (p *processor) subscribe(funnel string) (chan []byte, chan bool) {
	msgChan := make(chan []byte)
	stopChan := make(chan bool)

	go func() {
		for {
			p.transport.Subscribe(
				"#",
				p.name,
				funnel,
				msgChan,
				stopChan,
			)
		}
	}()

	return msgChan, stopChan
}
