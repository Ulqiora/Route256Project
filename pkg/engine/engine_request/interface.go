package engine_request

import (
	"context"

	"github.com/Ulqiora/Route256Project/pkg/engine/body"
)

type RequestEngine interface {
	GetMethodType() string
	GetPath() Path
	GetBody() body.Body
	Get(key string) string
	Has(key string) bool
	Context() context.Context
}

type Reader interface {
	Read() <-chan RequestEngine
}
type Writer interface {
	Write(str RequestEngine)
	Close()
}
