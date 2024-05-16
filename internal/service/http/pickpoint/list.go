package pickpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Ulqiora/Route256Project/internal/service/shared"
)

func (s *Service) List(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("list", req, req.Body)
	if s.cacher != nil {
		if bytes, err := s.cacher.Get(req.Context(), shared.HashRequest(req)); err != nil {
			err := json.Unmarshal(bytes, w)
			if err != nil {
				slog.Info(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
	data, err := s.controller.List(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	shared.WriteToWriter(w, data)
	if s.cacher != nil {
		if err := s.cacher.Set(req.Context(), shared.HashRequest(req), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
