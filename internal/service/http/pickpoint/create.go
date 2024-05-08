package pickpoint

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"homework/internal/model"
)

func (s *Service) Create(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	s.sendInfo("create", req, bytes.NewReader(body))
	var unm model.PickPoint
	if err := unm.Load(bytes.NewReader(body)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := s.controller.Create(req.Context(), unm)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("id", strconv.Itoa(int(id)))
	w.WriteHeader(http.StatusCreated)
}
