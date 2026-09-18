package errors

type ErrorResponse struct {
	ErrorValue APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

var ErrUnauthorized = APIError{Code: "Unauthorized", Message: "unauthorized"}
