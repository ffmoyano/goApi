package util

import (
	"encoding/json"
	"github.com/ffmoyano/goApi/logger"
	"net/http"
)

// Response sets headers and sends the response with the result of the dbPool query as a json
func Response(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err = w.Write(response); err != nil {
		logger.ErrorLogger.Printf("Couldn't write response:  %s", err)
	}
}
