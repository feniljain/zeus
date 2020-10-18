package views

import (
	"encoding/json"
	"net/http"
)

//Response depicts reponse structure to be sent
type Response struct {
	Code int `json:"code"`
	Meta `json:"meta"`
}

//Meta represents meta data to be sent
type Meta struct {
	Err     string                 `json:"err"`
	Payload map[string]interface{} `json:"payload"`
}

//SendResponse builds and sends the response
func SendResponse(w http.ResponseWriter, code int, err string, message map[string]interface{}) {
	resp := Response{
		Code: code,
		Meta: Meta{
			Err:     err,
			Payload: message,
		},
	}

	respJSON, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		message["message"] = "An error occurred"
		resp = Response{
			Code: code,
			Meta: Meta{
				Err:     err,
				Payload: message,
			},
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(respJSON)
}
