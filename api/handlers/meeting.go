package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	views "github.com/feniljain/zeus/api/views"
	pkg "github.com/feniljain/zeus/pkg"
	meeting "github.com/feniljain/zeus/pkg/meeting"
)

func createMeeting(meetingSvc meeting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := make(map[string]interface{})

		req := meeting.CreateMeetingReq{}

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&req)

		if err != nil {
			message["message"] = "An error occurred"
			views.SendResponse(w, http.StatusBadRequest, pkg.ErrWrongFormat.Error(), message)
			return
		}

		fmt.Println(req.Title)
		meetingSvc.CreateMeeting(req)

		message["message"] = "Successfully created meeting!"
		views.SendResponse(w, http.StatusCreated, "", message)
	}
}

//MakeParticipantHandler ...
func MakeMeetingHandler(meetingSvc meeting.Service) {
	http.HandleFunc("/meetings", createMeeting(meetingSvc))
}
