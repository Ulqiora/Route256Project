package engine_response

type ChanReader struct {
	source <-chan ResponseEngine
}

func (c *ChanReader) Read() <-chan ResponseEngine {
	return c.source
}

type ChanWriter struct {
	source chan<- ResponseEngine
}

func (c *ChanWriter) Write(response ResponseEngine) {
	c.source <- response
}

func (c *ChanWriter) Close() {
	close(c.source)
}

func CreateRWObjects(chanObj chan ResponseEngine) (*ChanReader, *ChanWriter) {
	return &ChanReader{source: chanObj}, &ChanWriter{source: chanObj}
}
