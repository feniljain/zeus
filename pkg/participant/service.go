package participant

//Service ...
type Service interface {
	CreateParticipant(CreateParticipantReq) error
}

type participantSvc struct {
	repo Repository
}

//MakeNewParticipantService ...
func MakeNewParticipantService(repo Repository) Service {
	return &participantSvc{repo: repo}
}

func (pSvc *participantSvc) CreateParticipant(req CreateParticipantReq) error {
	return pSvc.repo.CreateParticipant(req)
}
