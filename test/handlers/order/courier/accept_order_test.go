//go:build integration
// +build integration

package courier

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ulqiora/Route256Project/internal/core/http_service"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/model/order_changers"
	order_storage "github.com/Ulqiora/Route256Project/internal/repository/order/cli"
	pick_point_storage "github.com/Ulqiora/Route256Project/internal/repository/pick_point"
	pickpoint_storage "github.com/Ulqiora/Route256Project/internal/repository/pick_point/cli"
	"github.com/Ulqiora/Route256Project/internal/service/order/http/courier"
	"github.com/Ulqiora/Route256Project/test/postgresql"
	"github.com/stretchr/testify/require"
)

func TestAcceptOrderIntegrate(t *testing.T) {
	var (
		err      error = nil
		ctx            = context.Background()
		changers       = map[model.TypePacking]order_changers.ChangerOrder{
			model.TypePackage: order_changers.ChangerOrderPackage{},
			model.TypeBox:     order_changers.ChangerOrderBox{},
			model.TypeTape:    order_changers.ChangerOrderTape{},
		}
	)
	configfile := http_service.LoadConfig("/home/ccnc/GolandProjects/homework/config/config-testing.example.yaml")
	database, err := postgresql.NewTestDatabase(ctx, configfile)
	defer func() {
		fmt.Println(database.TruncateTable(ctx, "\"order\""))
		fmt.Println(database.TruncateTable(ctx, "pickpoint"))
	}()
	require.Nil(t, err)
	pickpoint_repo := pickpoint_storage.New(*database)
	id, err := pickpoint_repo.Create(ctx, pick_point_storage.PickPointDTO{
		Name:           "fds",
		Address:        "fds",
		ContactDetails: nil,
	})
	require.Nil(t, err)
	t.Run("accept order", func(t *testing.T) {
		repository := order_storage.New(*database)
		service := courier.New(repository, changers)
		handler := http.HandlerFunc(service.AcceptOrder)
		rec := httptest.NewRecorder()
		body := fmt.Sprintf(`{
            "customer_id": 1,
            "pick_point_id": %d,
            "shelf_life": "2024-04-22 20:52:51.961548 +00:00",
            "price": 12300,
            "weight": 9,
            "type": "Box"
        }`, id)
		req := httptest.NewRequest("POST", "/order", bytes.NewReader([]byte(body)))
		handler.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)
	})
	t.Run("err", func(t *testing.T) {
		t.Run("error date", func(t *testing.T) {
			repository := order_storage.New(*database)
			service := courier.New(repository, changers)
			handler := http.HandlerFunc(service.AcceptOrder)
			rec := httptest.NewRecorder()
			body := fmt.Sprintf(`{
                "customer_id": 1,
                "pick_point_id": %d,
                "shelf_life": "2024-04-31 20:52:51.961548 +00:00",
                "price": 12300,
                "weight": 9,
                "type": "Box"
            }`, id)
			req := httptest.NewRequest("POST", "/order", bytes.NewReader([]byte(body)))
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
		t.Run("error type", func(t *testing.T) {
			repository := order_storage.New(*database)
			service := courier.New(repository, changers)
			handler := http.HandlerFunc(service.AcceptOrder)
			rec := httptest.NewRecorder()
			body := fmt.Sprintf(`{
                "customer_id": 1,
                "pick_point_id": %d,
                "shelf_life": "2024-04-22 20:52:51.961548 +00:00",
                "price": 12300,
                "weight": 9,
                "type": "Box1"
            }`, id)
			req := httptest.NewRequest("POST", "/order", bytes.NewReader([]byte(body)))
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
		t.Run("error weight", func(t *testing.T) {
			repository := order_storage.New(*database)
			service := courier.New(repository, changers)
			handler := http.HandlerFunc(service.AcceptOrder)
			rec := httptest.NewRecorder()
			body := fmt.Sprintf(`{
                "customer_id": 1,
                "pick_point_id": %d,
                "shelf_life": "2024-04-22 20:52:51.961548 +00:00",
                "price": 12300,
                "weight": 100,
                "type": "Box1"
            }`, id)
			req := httptest.NewRequest("POST", "/order", bytes.NewReader([]byte(body)))
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

}
