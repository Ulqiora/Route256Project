package body

import "bytes"

type BodyType struct {
	data *bytes.Reader
}

func NewBody(allBytes []byte) *BodyType {
	return &BodyType{
		data: bytes.NewReader(allBytes),
	}
}

func (B *BodyType) Read(b []byte) (n int, err error) {
	return B.data.Read(b)
}
