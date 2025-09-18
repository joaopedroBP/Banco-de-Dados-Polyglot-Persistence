package main

import(
	"context"
	"fmt"
	"log"
	"os"
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main(){
	clientOptions := options.Client().ApplyURI("mongoDB://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.second)

	defer cancel()
	client,err := mongo.Connect(ctx,clientOptions)
	if err != nil {
		log.fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Conectado com sucesso!")

	
}

