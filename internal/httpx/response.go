package httpx

type Envelope struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(data any) Envelope {
	return Envelope{Code: 0, Message: "ok", Data: data}
}

func Error(code string, message string) Envelope {
	return Envelope{Code: 1, Message: message, Error: &ErrorBody{Code: code, Message: message}}
}
