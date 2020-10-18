package handlers

import (
	"fmt"
	"net/http"

	participant "github.com/feniljain/zeus/pkg/participant"
)

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Namastey Duniyaa!")
}

//MakeParticipantHandler ...
func MakeParticipantHandler(participantSvc participant.Service) {
	http.HandleFunc("/health", health)
}
