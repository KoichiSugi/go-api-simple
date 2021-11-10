package errs

type StatusCode int

const (
	OK                 StatusCode = 200
	Created            StatusCode = 201
	BadRequest         StatusCode = 400
	ServiceUnavailable StatusCode = 503
	Unauthorized       StatusCode = 401

	NotFound StatusCode = 404

	Internal  StatusCode = 500
	Forbidden StatusCode = 403
)
