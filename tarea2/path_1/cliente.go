package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	/*hello "example.com/m/tarea2/proto"*/
	ventas "example.com/m/tarea2/proto_ventas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Obtiene al momento de ejecutar el archivo el nombre del archivo JSON
	jsonFile := os.Args[1]

	// Lee el archivo JSON
	file, _ := os.ReadFile(jsonFile)

	// Crea una instancia de Order (estructura de Go)
	order := &ventas.OrderRequest{}

	// Convierte el archivo JSON a una estructura de Go (Order)
	_ = json.Unmarshal([]byte(file), order)

	// Se crea la coneccion con el servidor grpc
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := ventas.NewOrderProcessingClient(conn)

	respOrderProcess, err := c.ProcessOrder(context.Background(), &ventas.OrderRequest{
		Products: order.Products,
		Customer: order.Customer,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(respOrderProcess.Message)

}
