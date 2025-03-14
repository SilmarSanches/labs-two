package entities

type CustomErrors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CustomErrors) Error() string {
	return e.Message
}