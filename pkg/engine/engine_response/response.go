package engine_response

import "homework/pkg/engine/body"

const (
	StatusSuccess = iota
	StatusFail    = iota << 1
)

type ResponseRWType struct {
	body       body.Body
	statusCode int
	error      error
}

func New() *ResponseRWType {
	return &ResponseRWType{}
}

func (r *ResponseRWType) StatusCode() (int, error) {
	return r.statusCode, r.error
}

func (r *ResponseRWType) GetBody() body.Body {
	return r.body
}

func (r *ResponseRWType) SetBody(allBytes []byte) {
	r.body = body.NewBody(allBytes)
}
func (r *ResponseRWType) SetStatusCode(status int, err error) {
	r.statusCode = status
	r.error = err
}
