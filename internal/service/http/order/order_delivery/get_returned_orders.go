package order_delivery

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Ulqiora/Route256Project/internal/service/shared"
)

func (s *Service) GetReturnedOrders(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("get_returned_orders", req, req.Body)
	if s.cacher != nil {
		bytes, err := s.cacher.Get(req.Context(), shared.HashRequest(req))
		if err == nil {
			err := json.Unmarshal(bytes, w)
			if err != nil {
				slog.Info(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		slog.Info(err.Error())
	}
	values := req.URL.Query()
	orders, err := s.controller.GetReturnedOrders(req.Context(), &values)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shared.WriteToWriter(w, orders)
	if s.cacher != nil {
		if err := s.cacher.Set(req.Context(), shared.HashRequest(req), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
