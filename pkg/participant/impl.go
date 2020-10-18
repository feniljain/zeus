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
}
