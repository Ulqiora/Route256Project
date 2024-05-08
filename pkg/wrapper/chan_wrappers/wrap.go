package chan_wrappers

type ChanReader struct {
	source <-chan Request
}

func (c *ChanReader) Read() <-chan Request {
	return c.source
}

type ChanWriter struct {
	source chan<- Request
}

func (c *ChanWriter) Write(str Request) {
	c.source <- str
}

func (c *ChanWriter) Close() {
	close(c.source)
}

func CreateRWObjects(chanObj chan Request) (*ChanReader, *ChanWriter) {
	return &ChanReader{source: chanObj}, &ChanWriter{source: chanObj}
}
