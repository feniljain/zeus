package participant

//Repository ...
type Repository interface {
	CreateParticipant(CreateParticipantReq) error
}

//CreateParticipantReq represents the structure for creating a participant
type CreateParticipantReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Rsvp  string `json:"rsvp"`
}
