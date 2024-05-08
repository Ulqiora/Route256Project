package body

type Body interface {
	Read(b []byte) (n int, err error)
}
