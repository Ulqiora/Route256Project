package order_courier

import (
	"io"
	"net/http"

	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"github.com/gorilla/mux"
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
