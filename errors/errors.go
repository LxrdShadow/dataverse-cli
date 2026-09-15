package errors

type ErrorResponse struct {
	ErrorValue APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorValue.Error()
}

func (e APIError) Error() string {
	return e.Message
}
