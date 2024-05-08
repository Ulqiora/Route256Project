package engine

import (
	"homework/pkg/engine/body"
	"homework/pkg/engine/engine_response"
)

type ResponseReader struct {
	readerChan engine_response.Reader
}

type Response interface {
	GetBody() body.Body
	StatusCode() (int, error)
}

func NewResponseReader(reader engine_response.Reader) *ResponseReader {
	return &ResponseReader{
		readerChan: reader,
	}
}

func (reader *ResponseReader) Read() Response {
	return <-reader.readerChan.Read()
}
