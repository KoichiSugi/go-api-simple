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

func WrapError(info string, sc StatusCode, msg string) *RequestError {
	return &RequestError{
		Context: info,
		Code:    sc,
		Message: msg}
}
