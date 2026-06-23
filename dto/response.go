package dto

import "time"

type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      string    `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewErrorResponse(message, code string) ErrorResponse {
	return ErrorResponse{
		Error:     message,
		Code:      code,
		Timestamp: time.Now(),
	}
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Data:      data,
		Timestamp: time.Now(),
	}
}
