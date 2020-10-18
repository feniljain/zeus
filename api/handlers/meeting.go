package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			fmt.Println(err)
			message["message"] = "An error occurred"
			views.SendResponse(w, http.StatusBadRequest, pkg.ErrWrongFormat.Error(), message)
			return
		}

		err = meetingSvc.CreateMeeting(req)
		if err != nil {
			message["message"] = err.Error()
			views.SendResponse(w, http.StatusBadRequest, "An error occurred", message)
			return
		}

		message["message"] = "Successfully created meeting!"
		views.SendResponse(w, http.StatusCreated, "", message)
	}
}

func getMeetingDetails(meetingSvc meeting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := make(map[string]interface{})

		uri := r.URL.RequestURI()
		idx := strings.LastIndex(uri, "/")
		uid, err := strconv.Atoi(uri[idx+1:])
		if err != nil {
			message["message"] = "An error occurred"
			views.SendResponse(w, http.StatusBadRequest, pkg.ErrWrongFormat.Error(), message)
			return
		}

		res, err := meetingSvc.GetMeetingDetails(uid)
		if err != nil {
			message["message"] = err.Error()
			views.SendResponse(w, http.StatusNotFound, "An error occurred", message)
			return
		}

		message["meeting"] = res
		views.SendResponse(w, http.StatusOK, "", message)
	}
}

//MakeParticipantHandler ...
func MakeMeetingHandler(meetingSvc meeting.Service) {
	http.HandleFunc("/meetings", createMeeting(meetingSvc))
	http.HandleFunc("/meetings/", getMeetingDetails(meetingSvc))
}
