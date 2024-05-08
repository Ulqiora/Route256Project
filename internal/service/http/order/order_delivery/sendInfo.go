package order_delivery

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"homework/internal/service/broker_io"
)

func (s *Service) sendInfo(methodName string, r *http.Request, body io.Reader) {
	if s.sender == nil {
		return
	}
	content, _ := io.ReadAll(body)
	urlPath := r.URL.Path
	vars := mux.Vars(r)
	s.sender.SendMessage(broker_io.RequestMessage{
		Url:        urlPath,
		Content:    string(content),
		MethodName: methodName,
		Headers:    vars,
	})
}
