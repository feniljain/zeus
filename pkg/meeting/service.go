package meeting

import (
	entities "github.com/feniljain/zeus/pkg/entities"
)

//Service ...
type Service interface {
	CreateMeeting(CreateMeetingReq) error
	GetMeetingDetailsFromId(uid int) (entities.Meeting, error)
	GetMeetingDetailsFromTime(start, end string) ([]entities.Meeting, error)
}

type meetingSvc struct {
	repo Repository
}

//MakeNewMeetingService ...
func MakeNewMeetingService(repo Repository) Service {
	return &meetingSvc{repo: repo}
}

func (mSvc *meetingSvc) CreateMeeting(req CreateMeetingReq) error {
	return mSvc.repo.CreateMeeting(req)
}

func (mSvc *meetingSvc) GetMeetingDetailsFromId(uid int) (entities.Meeting, error) {
	return mSvc.repo.GetMeetingDetailsFromId(uid)
}

func (mSvc *meetingSvc) GetMeetingDetailsFromTime(start, end string) ([]entities.Meeting, error) {
	return mSvc.repo.GetMeetingDetailsFromTime(start, end)
}

//import (
//	entities "github.com/feniljain/zeus/pkg/entities"
//)

//Repository ...
type Repository interface {
	CreateMeeting(CreateMeetingReq) error
	GetMeetingDetailsFromId(uid int) (entities.Meeting, error)
	GetMeetingDetailsFromTime(start, end string) ([]entities.Meeting, error)
}

//CreateMeetingReq represents the structure for creating a meeting
type CreateMeetingReq struct {
	Title        string                 `json:"title"`
	StartTime    string                 `json:"starttime"`
	EndTime      string                 `json:"endtime"`
	Participants []entities.Participant `json:"participants"`
}
