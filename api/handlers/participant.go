package handlers

import (
	"encoding/json"
	"net/http"

	views "github.com/feniljain/zeus/api/views"
	pkg "github.com/feniljain/zeus/pkg"
	participant "github.com/feniljain/zeus/pkg/participant"
)

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := make(map[string]interface{})

		message["message"] = "Successfully created entity!"
		views.SendResponse(w, http.StatusOK, "", message)
	}
}

func createParticipant(participantSvc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := make(map[string]interface{})

		req := participant.CreateParticipantReq{}

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&req)

		if err != nil {
			message["message"] = "An error occurred"
			views.SendResponse(w, http.StatusBadRequest, pkg.ErrWrongFormat.Error(), message)
			return
		}

		message["message"] = "Successfully created participant!"
		views.SendResponse(w, http.StatusCreated, "", message)
	}
}

//MakeParticipantHandler ...
func MakeParticipantHandler(participantSvc participant.Service) {
	http.HandleFunc("/health", health())
	http.HandleFunc("/create-participant", createParticipant(participantSvc))
}
