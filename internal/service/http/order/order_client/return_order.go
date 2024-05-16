package order_client

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/gorilla/mux"
)

func (s *Service) ReturnOrder(w http.ResponseWriter, req *http.Request) {
	s.sendInfo("return_order", req, req.Body)

	id, customer_id, status, err := s.validateArgsReturnOrder(mux.Vars(req))
	if err != nil {
		slog.Warn(err.Error())
		w.WriteHeader(status)
		return
	}
	status, err = s.callReturnOrder(req.Context(), id, customer_id)
	if err != nil {
		slog.Warn(err.Error())
	}
	w.WriteHeader(status)
}

func (s *Service) callReturnOrder(ctx context.Context, id uint64, customer_id uint64) (int, error) {
	err := s.controller.ReturnOrder(ctx, id, customer_id)
	if err != nil {
		if errors.Is(err, repository.ErrorOrderNotFounded) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *Service) validateArgsReturnOrder(URLParams map[string]string) (uint64, uint64, int, error) {
	idStr, ok := URLParams["id"]
	if !ok {
		return 0, 0, http.StatusBadRequest, errors.New("param id is not setted")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	customeridStr, ok := URLParams["customer_id"]
	if !ok {
		return 0, 0, http.StatusBadRequest, errors.New("param id is not setted")
	}

	customerID, err := strconv.ParseUint(customeridStr, 10, 64)
	if err != nil {
		return 0, 0, http.StatusBadRequest, errors.New("param id is not correct")
	}
	return id, customerID, http.StatusOK, nil
}
