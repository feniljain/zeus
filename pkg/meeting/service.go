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

func (pSvc *meetingSvc) CreateMeeting(req CreateParticipantReq) error {
	return pSvc.repo.CreateParticipant(req)
}
