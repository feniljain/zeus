package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	http.HandleFunc("/health", hello)

	if err := http.ListenAndServe(":8010", nil); err != nil {
		panic(err)
	}

	fmt.Println("Listening on port 8010")
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Namastey Duniyaa!")
}

func initDatabase(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testing").Collection("participants")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.M{"name": "someone", "email": "someone@gmail.com", "rsvp": "yes"})
	id := res.InsertedID
	fmt.Println(id)
}
