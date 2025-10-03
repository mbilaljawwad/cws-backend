package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func setHeaders(writer http.ResponseWriter, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
}

func SendJSONResponse(writer http.ResponseWriter, statusCode int, data interface{}) {
	setHeaders(writer, statusCode)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SendJSONError(writer http.ResponseWriter, statusCode int, message string) {
	setHeaders(writer, statusCode)
	errResponse := ErrorResponse{
		Message:    message,
		StatusCode: statusCode,
	}

	SendJSONResponse(writer, statusCode, errResponse)
}
