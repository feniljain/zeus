package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/feniljain/zeus/api/handlers"
	meeting "github.com/feniljain/zeus/pkg/meeting"
	participant "github.com/feniljain/zeus/pkg/participant"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initialize() (meeting.Service, participant.Service) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	//Creating new instance of repositories
	participantRepo := participant.MakeNewParticipantRepo(client)
	meetingRepo := meeting.MakeNewMeetingRepo(client)

	//Injecting them in services and creating their new instances
	participantSvc := participant.MakeNewParticipantService(participantRepo)
	meetingSvc := meeting.MakeNewMeetingService(meetingRepo)

	//Injecting services and defining all handlers
	handlers.MakeParticipantHandler(participantSvc)
	handlers.MakeMeetingHandler(meetingSvc, participantSvc)

	return meetingSvc, participantSvc
}

func TestGetMeetingDetailsfromID(t *testing.T) {

	//Initial Setup
	meetingSvc, participantSvc := initialize()

	//Testing GetMeetingDetailsFromID
	req, err := http.NewRequest("GET", "/meetings/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetMeetingDetailsFromID(meetingSvc))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/meetings/4", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.GetMeetingDetailsFromID(meetingSvc))
	handler.ServeHTTP(rr, req)
	//fmt.Println(rr.Code)
	//fmt.Println(http.StatusNotFound)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Testing get all meetings of a participant
	req, err = http.NewRequest("GET", "/meetings?participant=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.BaseSubPathResolver(meetingSvc, participantSvc))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Testing get all meetings within a timeframe
	req, err = http.NewRequest("GET", "/meetings?start=2006-01-02 15:04:05&end=2007-01-02 15:34:05", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.BaseSubPathResolver(meetingSvc, participantSvc))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//Testing create a new meeting
	jsonStr := []byte(`{
    "title": "s",
    "starttime": "2006-01-02 15:04:05",
    "endtime": "2006-01-02 15:34:05",
    "participants": [
        {
            "name": "s",
            "email": "s",
            "rsvp": "Yes"
        },
        {
            "name": "se",
            "email": "sw",
            "rsvp": "Yes"
        }
    ]
}`)
	req, err = http.NewRequest("POST", "/meetings", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.BaseSubPathResolver(meetingSvc, participantSvc))
	handler.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
