package errorhandling

type StatusCode int

const (
	OK StatusCode = 200

	BadRequest         StatusCode = 400
	ServiceUnavailable StatusCode = 503
	Unauthorized       StatusCode = 401

	NotFound StatusCode = 404

	Internal  StatusCode = 500
	Forbidden StatusCode = 403
)
