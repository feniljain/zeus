package handlers

import (
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Namastey Duniyaa!")
}

//MakeParticipantHandler ...
func MakeParticipantHandler() {
	http.HandleFunc("/health", health)
}
