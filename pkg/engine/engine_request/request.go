package engine_request

import (
	"context"

	"github.com/Ulqiora/Route256Project/pkg/engine/body"
)

type MethodType string
type Path string

const (
	METHOD_GET  = "GET"
	METHOD_POST = "POST"
)

type RequestType struct {
	TypeMethod MethodType
	Body       body.Body
	Header     Header
	URL        Path
	ctx        context.Context
}

func (r *RequestType) GetBody() body.Body {
	return r.Body
}

func (r *RequestType) GetPath() Path {
	return r.URL
}

func (r *RequestType) Context() context.Context {
	return r.ctx
}

func (r *RequestType) Get(key string) string {
	return r.Header.content[key]
}
func (r *RequestType) Has(key string) bool {
	_, ok := r.Header.content[key]
	return ok
}

func (r *RequestType) GetMethodType() string {
	return string(r.TypeMethod)
}
