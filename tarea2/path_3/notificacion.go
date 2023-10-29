package main

import (
	"encoding/json"
	"log"

	ventas "example.com/m/tarea2/proto_ventas"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Request    *ventas.OrderRequest
	InsertedID string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Product struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Genre       string  `json:"genre"`
	Pages       int     `json:"pages"`
	Publication string  `json:"publication"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// Define the location struct
type Location struct {
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// Define the customer struct
type Customer struct {
	Name     string   `json:"name"`
	LastName string   `json:"lastname"`
	Email    string   `json:"email"`
	Location Location `json:"location"`
	Phone    string   `json:"phone"`
}

type RequestData struct {
	OrderID  string            `json:"orderID"`
	GroupID  string            `json:"groupID"`
	Products []*ventas.Product `json:"products"`
	Customer *ventas.Client    `json:"customer"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Deserializa el JSON en una instancia de Message que contiene req e insertedID
			var message Message
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Error al deserializar el mensaje: %v", err)
			} else {
				log.Printf(" [x] Received message: %v", message.Request.Customer.Name)
				log.Printf(" [x] Received insertedID: %s", message.InsertedID)
				log.Printf(" [x] El email esta siendo enviado...")
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
