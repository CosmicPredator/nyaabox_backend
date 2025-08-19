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