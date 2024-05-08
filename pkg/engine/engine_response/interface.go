package engine_response

import "homework/pkg/engine/body"

type ResponseEngine interface {
	GetBody() body.Body
	StatusCode() (int, error)
	SetBody(allBytes []byte)
	SetStatusCode(status int, err error)
}

type Reader interface {
	Read() <-chan ResponseEngine
}
type Writer interface {
	Write(response ResponseEngine)
	Close()
}
