package pickpoint

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Service) Delete(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("delete", req, req.Body)
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.controller.Delete(req.Context(), uint64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
