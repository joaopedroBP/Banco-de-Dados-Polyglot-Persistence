package main

import(
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main(){
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		log.Fatal(err)
	}
	
	err = client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal("falha ao conctar")
	}
	fmt.Println("conectado")
}

