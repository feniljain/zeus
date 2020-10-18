package handlers

import (
	"fmt"
	"net/http"

	views "github.com/feniljain/zeus/api/views"
	participant "github.com/feniljain/zeus/pkg/participant"
)

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Namastey Duniyaa!")
}

func createParticipant(participantSvc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		message := make(map[string]interface{})

		err := participantSvc.CreateParticipant(participant.CreateParticipantReq{
			Name:  "someone",
			Email: "someone@gmail.com",
			Rsvp:  "Yes",
		})

		if err != nil {
			message["message"] = "An error occurred"
			views.SendResponse(w, http.StatusInternalServerError, err.Error(), message)
			return
		}

		message["message"] = "Successfully created entity!"
		views.SendResponse(w, http.StatusCreated, "", message)
	}
}

//MakeParticipantHandler ...
func MakeParticipantHandler(participantSvc participant.Service) {
	http.HandleFunc("/health", health)
	http.HandleFunc("/create-participant", createParticipant(participantSvc))
}
