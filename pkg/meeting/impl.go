package meeting

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/feniljain/zeus/pkg"
	"github.com/feniljain/zeus/pkg/entities"
	_participant "github.com/feniljain/zeus/pkg/participant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
}

//MakeNewMeetingRepo takes and instance of mongo client and initializes the repo
func MakeNewMeetingRepo(client *mongo.Client) Repository {
	return &repo{client: client}
}

func (r *repo) GetMeetingDetails(uid int) (entities.Meeting, error) {
	participantCollection := r.client.Database("testing").Collection("meetings")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entities.Meeting

	filter := bson.M{"uid": uid}
	err := participantCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return entities.Meeting{}, pkg.ErrInvalidMeetingID
	}

	return res, nil
}

//CreateMeeting ...
func (r *repo) CreateMeeting(req CreateMeetingReq) error {
	fmt.Println("Create New Meeting")

	layout := "2006-01-02 15:04:05"
	//t, err := time.Parse(layout, "2014-11-17 23:02:03")
	//if err != nil {
	//	return pkg.ErrWrongTimestampFormat
	//}
	//y, m, d := t.Date()
	//fmt.Println(d, m, y, err)
	//h, mi, s := t.Clock()

	meeting := req
	meetingCollection := r.client.Database("testing").Collection("meetings")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := time.Parse(layout, meeting.StartTime)
	if err != nil {
		return pkg.ErrWrongTimestampFormat
	}

	_, err = time.Parse(layout, meeting.EndTime)
	if err != nil {
		return pkg.ErrWrongTimestampFormat
	}

	rand.Seed(time.Now().UTC().UnixNano())
	uid := rand.Intn(10000000)

	result, err := meetingCollection.InsertOne(ctx, bson.M{"uid": uid, "title": meeting.Title, "starttime": meeting.StartTime, "endtime": meeting.EndTime, "creationTimestamp": time.Now().Format(layout)})
	if err != nil {
		return err
	}

	for i := 0; i < len(req.Participants); i++ {
		participant := req.Participants[i]

		participantCollection := r.client.Database("testing").Collection("participants")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res _participant.CreateParticipantReq

		filter := bson.M{"email": participant.Email}
		err := participantCollection.FindOne(ctx, filter).Decode(&res)
		if err != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			id := fmt.Sprintf("%v", result.InsertedID)
			var meetingIDs = []string{id}

			_, err := participantCollection.InsertOne(ctx, bson.M{"name": participant.Name, "email": participant.Email, "rsvp": participant.Rsvp, "meetings": meetingIDs})
			if err != nil {
				return pkg.ErrInternalServer
			}

			//fmt.Println("Inserted participant")
			//fmt.Println(res.InsertedID)
		}
	}

	//fmt.Println("Inserted meeting")
	//fmt.Println(res.InsertedID)

	return nil
}
