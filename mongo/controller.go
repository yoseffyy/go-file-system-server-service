package mongo

import (
	"fmt"
	"log"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Controller struct {
	Client *mongo.Client
}

func NewController(url string) Controller {

	var newController Controller;

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	
	err = client.Connect(context.TODO())
    if err != nil {
        log.Fatal(err)
	}
	
	err = client.Ping(context.TODO(), nil);
	if err != nil {
		panic(err);
	}

	fmt.Println("Conncet to db");
	
	newController.Client = client;
	
	return newController;

}

