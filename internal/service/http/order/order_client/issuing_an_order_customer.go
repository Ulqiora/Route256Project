package order_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	OrderRespository "homework/internal/repository"
)

func (s *Service) IssuingAnOrderCustomer(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	s.sendInfo("issuing_an_customer_order", req, bytes.NewReader(body))
	ids, status, err := s.loadArgsIssuingAnOrderCustomer(bytes.NewReader(body))
	if err != nil {
		w.WriteHeader(status)
		return
	}
	status, err = s.callIssuingAnOrderCustomer(req.Context(), ids)
	w.WriteHeader(status)
}

func (s *Service) callIssuingAnOrderCustomer(ctx context.Context, ids []uint64) (int, error) {
	err := s.controller.IssuingToCustomer(ctx, ids)
	if err != nil {
		if errors.Is(err, OrderRespository.ErrorDataBase) {
			return http.StatusInternalServerError, err
		}
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *Service) loadArgsIssuingAnOrderCustomer(r io.Reader) ([]uint64, int, error) {
	byteAll, err := io.ReadAll(r)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var orderIDs []uint64
	if err = json.Unmarshal(byteAll, &orderIDs); err != nil {
		return nil, http.StatusBadRequest, err
	}
	return orderIDs, http.StatusOK, nil
}
