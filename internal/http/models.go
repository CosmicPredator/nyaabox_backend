package http

import "time"

type healthResponse struct {
	Status string    `json:"status"`
	Stamp  time.Time `json:"timestamp"`
}

func newHealthResponse(status string) *healthResponse {
	return &healthResponse{
		Status: status,
		Stamp:  time.Now(),
	}
}

type uploadSuccessResponse struct {
	healthResponse
	Files []string `json:"files"`
}

func newUploadSuccessResponse(files []string) *uploadSuccessResponse {
	return &uploadSuccessResponse{
		healthResponse: healthResponse {
			Status: "OK",
			Stamp: time.Now(),
		},
		Files: files,
	}
}

type ErrorResponse struct {
	healthResponse
	Message string `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		healthResponse: healthResponse{
			Status: "Not Ok",
			Stamp: time.Now(),
		},
		Message: message,
	}
}