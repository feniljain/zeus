package entities

//Participant struct represents participant structure in the database and json transactions
type Participant struct {
	name  string `json:"name"`
	email string `json:"email"`
	rsvp  string `json:"rsvp"`
}
