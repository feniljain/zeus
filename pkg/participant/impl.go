package participant

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
}

//MakeNewParticipantRepo takes and instance of mongo client and initializes the repo
func MakeNewParticipantRepo(client *mongo.Client) Repository {
	return &repo{client: client}
}

func (r *repo) CreateParticipant(req CreateParticipantReq) error {
	fmt.Println("Create New Participant")

	collection := r.client.Database("testing").Collection("participants")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, bson.M{"name": "someone", "email": "someone@gmail.com", "rsvp": "yes"})

	if err != nil {
		return err
	}

	id := res.InsertedID
	fmt.Println(id)

	return nil
}
