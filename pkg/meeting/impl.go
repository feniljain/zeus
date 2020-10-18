package meeting

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/feniljain/zeus/pkg"
	"github.com/feniljain/zeus/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	client *mongo.Client
}

var (
	layout string = "2006-01-02 15:04:05"
)

//MakeNewMeetingRepo takes and instance of mongo client and initializes the repo
func MakeNewMeetingRepo(client *mongo.Client) Repository {
	return &repo{client: client}
}

func (r *repo) GetMeetingDetailsFromTime(start, end string) ([]entities.Meeting, error) {
	startGiven, err := time.Parse(layout, start)
	if err != nil {
		return nil, pkg.ErrWrongTimestampFormat
	}

	endGiven, err := time.Parse(layout, end)
	if err != nil {
		return nil, pkg.ErrWrongTimestampFormat
	}

	var meetings []entities.Meeting

	meetingCollection := r.client.Database("testing").Collection("meetings")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"title": "s"}
	cur, err := meetingCollection.Find(ctx, filter)
	if err != nil {
		return nil, pkg.ErrInternalServer
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result entities.Meeting
		err := cur.Decode(&result)
		if err != nil {
			return nil, pkg.ErrInternalServer
		}

		starttime, err := time.Parse(layout, result.StartTime)
		if err != nil {
			return nil, pkg.ErrInternalServer
		}

		endtime, err := time.Parse(layout, result.EndTime)
		if err != nil {
			return nil, pkg.ErrInternalServer
		}

		if startGiven.After(starttime) || endGiven.After(endtime) || (startGiven.Equal(starttime) && endGiven.Equal(endtime)) {
			meetings = append(meetings, result)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, pkg.ErrInternalServer
	}

	return meetings, nil
}

func (r *repo) GetMeetingDetailsFromId(uid int) (entities.Meeting, error) {
	meetingCollection := r.client.Database("testing").Collection("meetings")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entities.Meeting

	filter := bson.M{"uid": uid}
	err := meetingCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return entities.Meeting{}, pkg.ErrInvalidMeetingID
	}

	return res, nil
}

//CreateMeeting ...
func (r *repo) CreateMeeting(req CreateMeetingReq) error {
	fmt.Println("Create New Meeting")

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

	_, err = meetingCollection.InsertOne(ctx, bson.M{"uid": uid, "title": meeting.Title, "starttime": meeting.StartTime, "endtime": meeting.EndTime, "creationTimestamp": time.Now().Format(layout)})
	if err != nil {
		return err
	}

	for i := 0; i < len(req.Participants); i++ {
		participant := req.Participants[i]

		participantCollection := r.client.Database("testing").Collection("participants")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res entities.Participant

		filter := bson.M{"email": participant.Email}
		err := participantCollection.FindOne(ctx, filter).Decode(&res)
		if err != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var meetingIDs = []int{uid}

			_, err := participantCollection.InsertOne(ctx, bson.M{"name": participant.Name, "email": participant.Email, "rsvp": participant.Rsvp, "meetings": meetingIDs})
			if err != nil {
				return pkg.ErrInternalServer
			}
			continue
		}

		res.Meetings = append(res.Meetings, uid)

		flag := 0

		for i := 0; i < len(res.Meetings); i++ {
			uid1 := res.Meetings[i]
			m, err := r.GetMeetingDetailsFromId(uid1)
			if err != nil {
				return err
			}

			startGiven, err := time.Parse(layout, meeting.StartTime)
			endGiven, err := time.Parse(layout, meeting.EndTime)
			starttime, err := time.Parse(layout, m.StartTime)
			endtime, err := time.Parse(layout, m.EndTime)

			if (m.UID != uid) && !(endGiven.Before(starttime) || startGiven.After(endtime) || endGiven.Equal(starttime) || startGiven.Equal(endtime)) {
				fmt.Println(startGiven, endGiven, starttime, endtime)
				fmt.Println(res)
				fmt.Println(m.Title)
				flag = 1
				break
			}
		}

		if flag == 1 {
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter = bson.M{"email": res.Email}
		_, err = participantCollection.ReplaceOne(ctx, filter, bson.M{"name": res.Name, "email": res.Email, "rsvp": res.Rsvp, "meetings": res.Meetings})
		if err != nil {
			return pkg.ErrInternalServer
		}
	}

	return nil
}
