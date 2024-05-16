package order_courier

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/repository"
)

func (s *Service) AcceptOrder(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	s.sendInfo("accept_order", req, bytes.NewReader(body))
	var initData model.OrderInitData
	if err := initData.LoadFromRequest(bytes.NewReader(body)); err != nil {
		slog.Warn(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, status, err := s.callAcceptOrder(req.Context(), initData)
	if err != nil {
		slog.Warn(err.Error())
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Id", strconv.Itoa(int(id)))
	w.WriteHeader(status)
}

func (s *Service) callAcceptOrder(ctx context.Context, initdata model.OrderInitData) (uint64, int, error) {

	id, err := s.controller.AcceptOrder(ctx, initdata)
	if err != nil {
		if errors.Is(err, repository.ErrorOrderNotCreated) {
			return 0, http.StatusNotFound, err
		}
		return 0, http.StatusInternalServerError, err
	}
	return id, http.StatusCreated, err
}
