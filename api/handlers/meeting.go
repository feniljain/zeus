package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	views "github.com/feniljain/zeus/api/views"
	pkg "github.com/feniljain/zeus/pkg"
	"github.com/feniljain/zeus/pkg/entities"
	meeting "github.com/feniljain/zeus/pkg/meeting"
	"github.com/feniljain/zeus/pkg/participant"
)

//Handle participant query
func participantQuery(participantSvc participant.Service, meetingSvc meeting.Service, w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	message := make(map[string]interface{})
	var meetings []entities.Meeting
	ids, err := participantSvc.GetMeetings(queries.Get("participant"))

	if err != nil {
		message["message"] = err.Error()
		views.SendResponse(w, http.StatusBadRequest, "An error occurred", message)
		return
	}

	for i := 0; i < len(ids); i++ {
		meeting, err := meetingSvc.GetMeetingDetailsFromId(ids[i])

		if err != nil {
			message["message"] = err.Error()
			views.SendResponse(w, http.StatusBadRequest, "An error occurred", message)
		}

		meetings = append(meetings, meeting)
	}

	message["meetings"] = meetings
	views.SendResponse(w, http.StatusOK, "", message)
	return
}

//Handle time query
func timeQuery(meetingSvc meeting.Service, w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	message := make(map[string]interface{})
	meetings, err := meetingSvc.GetMeetingDetailsFromTime(queries["start"][0], queries["end"][0])
	if err != nil {
		fmt.Println(err)
		message["message"] = pkg.ErrInternalServer.Error()
		views.SendResponse(w, http.StatusInternalServerError, "An error occurred", message)
		return
	}

	message["meetings"] = meetings
	views.SendResponse(w, http.StatusOK, "", message)
	return
}

//Handle create meeting
func createMeeting(meetingSvc meeting.Service, w http.ResponseWriter, r *http.Request) {
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
	return
}

//Resolve if there are queries, then redirect them to needed method or create a new meeting
func baseSubPathResolver(meetingSvc meeting.Service, participantSvc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queries := r.URL.Query()

		if queries.Get("participant") != "" {
			participantQuery(participantSvc, meetingSvc, w, r)
			return
		}

		if queries.Get("start") != "" {
			timeQuery(meetingSvc, w, r)
			return
		}

		createMeeting(meetingSvc, w, r)
		return
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

		res, err := meetingSvc.GetMeetingDetailsFromId(uid)
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
func MakeMeetingHandler(meetingSvc meeting.Service, participantSvc participant.Service) {
	http.HandleFunc("/meetings", baseSubPathResolver(meetingSvc, participantSvc))
	http.HandleFunc("/meetings/", getMeetingDetails(meetingSvc))
}
