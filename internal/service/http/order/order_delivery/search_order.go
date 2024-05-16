package order_delivery

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/service/shared"
	"github.com/gorilla/mux"
)

func (s *Service) SearchOrders(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("search_order", req, req.Body)
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
	vars := mux.Vars(req)
	idCustomer, err := strconv.ParseUint(vars["customer_id"], 10, 64)
	orders, err := s.controller.SearchOrders(req.Context(), idCustomer, req.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	shared.WriteToWriter(w, orders)
	if s.cacher != nil {
		if err := s.cacher.Set(req.Context(), shared.HashRequest(req), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
