package main

import (
	"log"

	"github.com/upfluence/sensu-handler-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/event"
	"github.com/upfluence/sensu-handler-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
	"github.com/upfluence/sensu-handler-go/sensu"
	"github.com/upfluence/sensu-handler-go/sensu/handler"
)

type mockHandler struct{}

func (h mockHandler) Handle(e *event.Event) error {
	log.Printf("%+v\n", e)
	return nil
}

func main() {
	handler.Store["testtest"] = mockHandler{}

	sensu.NewHandler(
		rabbitmq.NewRabbitMQTransport("amqp://localhost:5672/%2f"),
	).Start()
}
