package meeting

import (
	"context"
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

//GetMeetingDetailsFromTime returns the meetings which lie in specified timeframes
func (r *repo) GetMeetingDetailsFromTime(start, end string) ([]entities.Meeting, error) {
	//Check if the dates are consitent with the format
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

	//Get all Meeetings
	filter := bson.M{}
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

		//Extracting time of current meeting in consideration
		starttime, err := time.Parse(layout, result.StartTime)
		if err != nil {
			return nil, pkg.ErrInternalServer
		}

		endtime, err := time.Parse(layout, result.EndTime)
		if err != nil {
			return nil, pkg.ErrInternalServer
		}

		//Checking if current meeting lies in given timeframe
		if startGiven.After(starttime) || endGiven.After(endtime) || (startGiven.Equal(starttime) && endGiven.Equal(endtime)) {
			meetings = append(meetings, result)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, pkg.ErrInternalServer
	}

	return meetings, nil
}

//GetMeetingDetailsFromId queries for a meeting with the given ID
func (r *repo) GetMeetingDetailsFromId(uid int) (entities.Meeting, error) {
	meetingCollection := r.client.Database("testing").Collection("meetings")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res entities.Meeting

	//Finding meeting given id by setting up the filter and using FindOne
	filter := bson.M{"uid": uid}
	err := meetingCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return entities.Meeting{}, pkg.ErrInvalidMeetingID
	}

	return res, nil
}

//CreateMeeting is used to make a new meeting and also register new participants if they dont have a meeting at the same time
func (r *repo) CreateMeeting(req CreateMeetingReq) error {

	meeting := req
	meetingCollection := r.client.Database("testing").Collection("meetings")

	//Creating a randome UID
	rand.Seed(time.Now().UTC().UnixNano())
	uid := rand.Intn(10000000)

	//Creating new meeting
	err := createNewMeeting(meeting, meetingCollection, uid)
	if err != nil {
		return err
	}

	//Going through participants to register them or to decide whether to include them or not
	for i := 0; i < len(req.Participants); i++ {
		participant := req.Participants[i]

		participantCollection := r.client.Database("testing").Collection("participants")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res entities.Participant

		//Finding if the participant exists
		filter := bson.M{"email": participant.Email}
		err := participantCollection.FindOne(ctx, filter).Decode(&res)
		if err != nil {
			//Creating participant as it doesn't exist
			err := createNewParticipant(participantCollection, participant, uid)
			if err != nil {
				return err
			}
			continue
		}

		res.Meetings = append(res.Meetings, uid)

		//Checking is theres a time clash
		yes, err := isTimeClash(r, res.Meetings, meeting, uid)
		if err != nil {
			return err
		}

		if yes {
			continue
		}

		//Updating Participant Info if there isnt any clash
		err = updateParticipantInfo(res, participantCollection)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateParticipantInfo(participant entities.Participant, participantCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": participant.Email}
	_, err := participantCollection.ReplaceOne(ctx, filter, bson.M{"name": participant.Name, "email": participant.Email, "rsvp": participant.Rsvp, "meetings": participant.Meetings})
	if err != nil {
		return pkg.ErrInternalServer
	}

	return nil
}

func createNewMeeting(currMeeting CreateMeetingReq, meetingCollection *mongo.Collection, uid int) error {
	_, err := time.Parse(layout, currMeeting.StartTime)
	if err != nil {
		return pkg.ErrWrongTimestampFormat
	}

	_, err = time.Parse(layout, currMeeting.EndTime)
	if err != nil {
		return pkg.ErrWrongTimestampFormat
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = meetingCollection.InsertOne(ctx, bson.M{"uid": uid, "title": currMeeting.Title, "starttime": currMeeting.StartTime, "endtime": currMeeting.EndTime, "creationTimestamp": time.Now().Format(layout)})
	if err != nil {
		return err
	}

	return nil
}

func isTimeClash(r *repo, meetings []int, currMeeting CreateMeetingReq, uid int) (bool, error) {
	flag := 0

	for i := 0; i < len(meetings); i++ {
		uid1 := meetings[i]
		m, err := r.GetMeetingDetailsFromId(uid1)
		if err != nil {
			return false, err
		}

		//Extracting all needed timeframes
		startGiven, err := time.Parse(layout, currMeeting.StartTime)
		endGiven, err := time.Parse(layout, currMeeting.EndTime)
		starttime, err := time.Parse(layout, m.StartTime)
		endtime, err := time.Parse(layout, m.EndTime)

		//Checking for clash
		if (m.UID != uid) && !(endGiven.Before(starttime) || startGiven.After(endtime) || endGiven.Equal(starttime) || startGiven.Equal(endtime)) {
			//fmt.Println(startGiven, endGiven, starttime, endtime)
			//fmt.Println(res)
			//fmt.Println(m.Title)
			flag = 1
			break
		}
	}

	if flag == 1 {
		return true, nil
	}

	return false, nil
}

func createNewParticipant(participantCollection *mongo.Collection, participant entities.Participant, uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var meetingIDs = []int{uid}

	_, err := participantCollection.InsertOne(ctx, bson.M{"name": participant.Name, "email": participant.Email, "rsvp": participant.Rsvp, "meetings": meetingIDs})
	if err != nil {
		return pkg.ErrInternalServer
	}

	return nil
}
