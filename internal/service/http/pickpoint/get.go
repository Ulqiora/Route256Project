package pickpoint

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/Ulqiora/Route256Project/internal/service/shared"
	"github.com/gorilla/mux"
)

func (s *Service) Get(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("get", req, req.Body)
	key, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.Atoi(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dataDTO, err := s.controller.GetByID(req.Context(), uint64(keyInt))
	if err != nil {
		if errors.Is(err, repository.ErrorObjectNotFounded) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	shared.WriteToWriter(w, dataDTO)

}
