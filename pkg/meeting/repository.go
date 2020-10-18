package meeting

import (
	entities "github.com/feniljain/zeus/pkg/entities"
)

//Repository ...
type Repository interface {
	CreateMeeting(CreateMeetingReq) error
}

//CreateMeetingReq represents the structure for creating a meeting
type CreateMeetingReq struct {
	UID          string                 `json:"uid"`
	Title        string                 `json:"name"`
	StartTime    string                 `json:"starttime"`
	EndTime      string                 `json:"endtime"`
	Participants []entities.Participant `json:"participants"`
}
