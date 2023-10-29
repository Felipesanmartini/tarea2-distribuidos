package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	ventas "example.com/m/tarea2/proto_ventas"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Request    *ventas.OrderRequest
	InsertedID string
}

type Book struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Author      string             `bson:"author"`
	Genre       string             `bson:"genre"`
	Pages       int                `bson:"pages"`
	Publication string             `bson:"publication"`
	Quantity    int                `bson:"quantity"`
	Price       float64            `bson:"price"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
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

				// Configura la conexión a la base de datos MongoDB
				var ctx context.Context
				client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@tarea2.6awdbqv.mongodb.net/?retryWrites=true&w=majority"))
				if err != nil {
					return
				}
				defer client.Disconnect(ctx)

				// Accede a la colección
				// Accede a la base de datos y la colección
				database := client.Database("Tarea2")
				collection := database.Collection("products")

				for _, product := range message.Request.Products {
					// Define el filtro para encontrar el documento que deseas actualizar
					filter := bson.M{"title": product.Title}

					var book Book
					err = collection.FindOne(context.Background(), filter).Decode(&book)
					if err != nil {
						log.Fatal(err)
					}

					// Define la actualización que agregarás al documento
					update := bson.M{
						"$set": bson.M{"quantity": book.Quantity - int(product.Quantity)}, // Cambia "110" al valor deseado
					}

					// Realiza la actualización
					_, err = collection.UpdateOne(context.Background(), filter, update)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("Documento actualizado con éxito para el producto %s.\n", product.Title)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
