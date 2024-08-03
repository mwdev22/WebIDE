package handlers

func RaiseApiError(message string, code int) {

}

type ApiError interface {
	error
	Error() string
	StatusCode()
}

type BadRequest struct {
	Message string
	Code    int
}

func (e *BadRequest) Error() string {
	return e.Message
}

type Unathorized struct {
	Message    string
	StatusCode int
}

type NotFound struct {
	Message    string
	StatusCode int
}
