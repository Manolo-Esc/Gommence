package ports

// from https://www.joeshaw.org/error-handling-in-go-http-applications/

type APIError interface {
	// APIError returns an HTTP status code and an error message.
	APIError() (int, string)
	Error() string
	Status() int
}

type structAPIError struct {
	status int
	msg    string
}

func (e structAPIError) Error() string {
	return e.msg
}

func (e structAPIError) Status() int {
	return e.status
}

func (e structAPIError) APIError() (int, string) {
	return e.status, e.msg
}

func NewAPIError(status int, msg string) APIError {
	return &structAPIError{status, msg}
}
