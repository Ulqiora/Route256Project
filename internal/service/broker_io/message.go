package broker_io

type RequestMessage struct {
	Url        string
	Content    string
	MethodName string
	Headers    map[string]string
}
