package meeting

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
}

//MakeNewMeetingRepo takes and instance of mongo client and initializes the repo
func MakeNewMeetingRepo(client *mongo.Client) Repository {
	return &repo{client: client}
}

//CreateMeeting ...
func (r *repo) CreateMeeting(req CreateMeetingReq) error {
	fmt.Println("Create New Meeting")

	//uid := rand.Intn(10000000)

	rand.Seed(time.Now().UTC().UnixNano())

	//collection := r.client.Database("testing").Collection("meeting")

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	//res, err := collection.InsertOne(ctx, bson.M{"uid": uid, "name": "someone", "email": "someone@gmail.com", "rsvp": "yes"})

	//if err != nil {
	//	return err
	//}

	//id := res.InsertedID
	//fmt.Println(id)

	return nil
}
