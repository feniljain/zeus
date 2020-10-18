package participant

import (
	"context"
	"fmt"
	"time"

	"github.com/feniljain/zeus/pkg"
	"github.com/feniljain/zeus/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
}

//MakeNewParticipantRepo takes and instance of mongo client and initializes the repo
func MakeNewParticipantRepo(client *mongo.Client) Repository {
	return &repo{client: client}
}

func (r *repo) GetMeetings(email string) ([]int, error) {
	participantCollection := r.client.Database("testing").Collection("participants")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entities.Participant

	filter := bson.M{"email": email}
	err := participantCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		fmt.Println(err)
		return nil, pkg.ErrParticipantEmail
	}

	return res.Meetings, nil
}

func (r *repo) CreateParticipant(req CreateParticipantReq) error {
	fmt.Println("Create New Participant")

	collection := r.client.Database("testing").Collection("participants")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, bson.M{"name": "someone", "email": "someone@gmail.com", "rsvp": "yes"})
	if err != nil {
		return err
	}

	id := res.InsertedID
	fmt.Println(id)

	return nil
}
