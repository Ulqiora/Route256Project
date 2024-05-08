package chan_wrappers

type Body interface {
	Read(b []byte) (n int, err error)
}

type Request interface {
	GetMethodType() string
	GetBody() Body
	Get(key string) string
}

type Response interface {
	GetBody() Body
	Get(key string) string
}

type Reader interface {
	Read() <-chan Request
}
type Writer interface {
	Write(str Request)
	Close()
}
