package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	hello "example.com/m/tarea2/proto"
	ventas "example.com/m/tarea2/proto_ventas"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedGreeterServer
	ventas.UnimplementedOrderProcessingServer
}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (s *server) ProcessOrder(ctx context.Context, req *ventas.OrderRequest) (*ventas.OrderResponse, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@tarea2.6awdbqv.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// Accede a la base de datos y la colección
	database := client.Database("Tarea2")
	collection := database.Collection("orders")

	insertResult, err := collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	insertedID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	sendMessageRabbitmq(req, insertedID)

	return &ventas.OrderResponse{Message: "Order is in process...\nAn e-mail will be sent to you, " + req.Customer.Name + ".\n\nYour order id is: " + insertedID}, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func sendMessageRabbitmq(req *ventas.OrderRequest, insertedID string) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := struct {
		Request    *ventas.OrderRequest
		InsertedID string
	}{
		Request:    req,
		InsertedID: insertedID,
	}

	reqJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to serialize OrderRequest to JSON: %v", err)
		return
	}

	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json", // Usar application/json para el contenido
			Body:        reqJSON,            // Usar el JSON serializado como cuerpo
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", req)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &server{})
	ventas.RegisterOrderProcessingServer(s, &server{})
	fmt.Println("Port Listening on: 50051")
	s.Serve(lis)

}
