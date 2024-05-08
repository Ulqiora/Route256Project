package pickpoint

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"homework/internal/model"
)

func (s *Service) Update(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	s.sendInfo("update", req, bytes.NewReader(body))
	var unm model.PickPoint
	if err := unm.Load(bytes.NewReader(body)); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	unm.ID = id
	idReturned, err := s.controller.Update(req.Context(), unm)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("id", strconv.Itoa(int(idReturned)))
}
