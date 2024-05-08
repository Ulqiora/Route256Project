package engine_request

type Header struct {
	content map[string]string
}

func NewHeader() Header {
	return Header{content: make(map[string]string)}
}

func (h *Header) Get(key string) string {
	return h.content[key]
}
func (h *Header) Set(key, value string) {
	h.content[key] = value
}
