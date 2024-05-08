package engine_request

type ChanReader struct {
	source <-chan RequestEngine
}

func (c *ChanReader) Read() <-chan RequestEngine {
	return c.source
}

type ChanWriter struct {
	source chan<- RequestEngine
}

func (c *ChanWriter) Write(str RequestEngine) {
	c.source <- str
}

func (c *ChanWriter) Close() {
	close(c.source)
}

func CreateRWObjects(chanObj chan RequestEngine) (*ChanReader, *ChanWriter) {
	return &ChanReader{source: chanObj}, &ChanWriter{source: chanObj}
}
