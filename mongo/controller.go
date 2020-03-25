package mongo

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Controller struct {
	Client *mongo.Client
}

func NewController(url string) Controller {
	controller := Controller{}
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conncet", url)

	controller.Client = client

	return controller
}
