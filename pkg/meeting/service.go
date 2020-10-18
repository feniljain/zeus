package meeting

//Service ...
type Service interface {
	CreateMeeting(CreateMeetingReq) error
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

//import (
//	entities "github.com/feniljain/zeus/pkg/entities"
//)

//Repository ...
type Repository interface {
	CreateMeeting(CreateMeetingReq) error
}

//CreateMeetingReq represents the structure for creating a meeting
type CreateMeetingReq struct {
	UID          string   `json:"uid"`
	Title        string   `json:"title"`
	StartTime    string   `json:"starttime"`
	EndTime      string   `json:"endtime"`
	Participants []string `json:"participants"`
}
