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
	docs := "www.mongodb.com/docs/drivers/go/current/"
	uri := "mongodb://localhost:27017"
	
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err != nil {
		panic("erro")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic("erro")
		}
	}()

	err = client.Ping(context.TODO(),nil)
	if err != nil{
		panic("erro")
	}

	fmt.Println("Conexão com o mongo efetuada com exito!")

	
}

