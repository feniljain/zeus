package handlers

import (
	"encoding/json"
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

	if r.Method != "GET" {
		message["message"] = "Requested URL not found, please check method of request and URL"
		views.SendResponse(w, http.StatusNotFound, "An error occurred", message)
		return
	}

	var page int = -1

	if queries.Get("page") != "" {
		p, err := strconv.Atoi(queries.Get("page"))
		if err != nil {
			message["message"] = err.Error()
			views.SendResponse(w, http.StatusBadRequest, "An error occurred", message)
			return
		}
		page = p
	}

	ids, err := participantSvc.GetAllMeetings(queries.Get("participant"), page)

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

	if r.Method != "GET" {
		message["message"] = "Requested URL not found, please check method of request and URL"
		views.SendResponse(w, http.StatusNotFound, "An error occurred", message)
		return
	}

	if err != nil {
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

	if r.Method != "POST" {
		message["message"] = "Requested URL not found, please check method of request and URL"
		views.SendResponse(w, http.StatusNotFound, "An error occurred", message)
		return
	}

	req := meeting.CreateMeetingReq{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
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

//BaseSubPathResolver resolves if there are queries, then redirect them to needed method or create a new meeting
func BaseSubPathResolver(meetingSvc meeting.Service, participantSvc participant.Service) http.HandlerFunc {
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

//GetMeetingDetails gets the details of a meeting from the ID
func GetMeetingDetailsFromID(meetingSvc meeting.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := make(map[string]interface{})

		if r.Method != "GET" {
			message["message"] = "Requested URL not found, please check method of request and URL"
			views.SendResponse(w, http.StatusNotFound, "An error occurred", message)
			return
		}

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
	http.HandleFunc("/meetings", BaseSubPathResolver(meetingSvc, participantSvc))
	http.HandleFunc("/meetings/", GetMeetingDetailsFromID(meetingSvc))
}
