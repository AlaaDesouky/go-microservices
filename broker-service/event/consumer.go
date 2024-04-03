package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Services struct {
	Log string
}

type Consumer struct {
	conn *amqp.Connection
	queueName string
	services Services
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection, services *Services) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
		services: *services,
	}

	if err := consumer.setup(); err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}


func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for m := range messages {
			var payload Payload
			_ = json.Unmarshal(m.Body, &payload)
			go consumer.handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)

	<-forever

	return nil
}

func (consumer *Consumer) handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := consumer.logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	default:
		err := consumer.logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func (consumer *Consumer) logEvent(e Payload) error {
	jsonData, _ := json.MarshalIndent(e, "", "\t")

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/log", consumer.services.Log), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
		if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted{
		return err
	}

	return nil
}
