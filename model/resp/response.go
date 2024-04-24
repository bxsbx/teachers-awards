package resp

type Response struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	ErrStack interface{} `json:"error_stack,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
