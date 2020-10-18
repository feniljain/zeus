package entities

//Participant struct represents participant structure in the database and json transactions
type Participant struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Rsvp     string `json:"rsvp"`
	Meetings []int  `json:"meetings"`
}

//Meeting struct represents participant structure in the database and json transactions
type Meeting struct {
	UID               int    `json:"uid"`
	Title             string `json:"title"`
	StartTime         string `json:"starttime"`
	EndTime           string `json:"endtime"`
	CreationTimeStamp string `json:"creationTimestamp"`
}
