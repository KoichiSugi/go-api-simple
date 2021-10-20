package errorhandling

import (
	"fmt"
)

//creating custom error
type RequestError struct {
	Context string
	Code    StatusCode
	Message string
}

func (r *RequestError) Error() string {
	return fmt.Sprintf(
		"Request Error Context: %s , StatusCode: %v, Message: %s", r.Context, r.Code, r.Message)
}
func (r *RequestError) StatusServiceUnavailable() bool {
	return r.Code == ServiceUnavailable //503
}

func (r *RequestError) StatusInternalServerError() bool {
	return r.Code == Internal //500
}

func (r *RequestError) StatusNotFound() bool {
	return r.Code == NotFound //404
}
func (r *RequestError) StatusBadRequest() bool {
	return r.Code == BadRequest //400
}

func WrapError(info string, sc StatusCode, msg string) *RequestError {
	return &RequestError{
		Context: info,
		Code:    sc,
		Message: msg}
}
