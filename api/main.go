package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	handlers "github.com/feniljain/zeus/api/handlers"
	meeting "github.com/feniljain/zeus/pkg/meeting"
	participant "github.com/feniljain/zeus/pkg/participant"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	//Initialize DB
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

	//Dependency Injection
	inject(client)

	fmt.Println("Listening on port 8010")

	//Setup Listener of requests
	if err := http.ListenAndServe(":8010", nil); err != nil {
		panic(err)
	}
}

func inject(client *mongo.Client) {
	//Creating new instance of repositories
	participantRepo := participant.MakeNewParticipantRepo(client)
	meetingRepo := meeting.MakeNewMeetingRepo(client)

	//Injecting them in services and creating their new instances
	participantSvc := participant.MakeNewParticipantService(participantRepo)
	meetingSvc := meeting.MakeNewMeetingService(meetingRepo)

	//Injecting services and defining all handlers
	handlers.MakeParticipantHandler(participantSvc)
	handlers.MakeMeetingHandler(meetingSvc, participantSvc)
}
