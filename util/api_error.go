package util

type ApiError struct {
	Message string
	Code    int
}

func NewApiError(message string, code int) *ApiError {
	return &ApiError{
		Message: message,
		Code:    code,
	}
}

func (err *ApiError) Error() string {
	return err.Message
}
