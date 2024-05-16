package order_courier

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/gorilla/mux"
)

func (s *Service) ReturnOrderToCourier(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("return_order_to_courier", req, req.Body)
	id, status, err := s.validateArgsReturnOrderToCourier(mux.Vars(req))
	if err != nil {
		slog.Warn(err.Error())
		w.WriteHeader(status)
		return
	}
	status, err = s.callReturnOrderToCourier(req.Context(), id)
	if err != nil {
		slog.Warn(err.Error())
	}
	w.WriteHeader(status)
}

func (s *Service) callReturnOrderToCourier(ctx context.Context, id uint64) (int, error) {
	err := s.controller.ReturnOrderToCourier(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorOrderNotFounded) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *Service) validateArgsReturnOrderToCourier(URLParams map[string]string) (uint64, int, error) {
	idStr, ok := URLParams["id"]
	if !ok {
		return 0, http.StatusBadRequest, errors.New("param id is not setted")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest, errors.New("param id is not correct")
	}
	return id, http.StatusOK, nil
}
