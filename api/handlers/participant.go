package handlers

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Namastey Duniyaa!")
}

//MakeParticipantHandler ...
func MakeParticipantHandler() {
	http.HandleFunc("/health", hello)
}
