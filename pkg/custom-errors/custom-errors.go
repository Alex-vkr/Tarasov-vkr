package custom_errors

type ErrHttp struct {
	Code    int
	Message string
}

func (e *ErrHttp) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func New(code int, message string) *ErrHttp {
	return &ErrHttp{Code: code, Message: message}
}
