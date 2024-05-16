package courier

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ulqiora/Route256Project/internal/core/http_service"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/model/order_changers"
	order_storage "github.com/Ulqiora/Route256Project/internal/repository/order/cli"
	pick_point_storage "github.com/Ulqiora/Route256Project/internal/repository/pick_point"
	pickpoint_storage "github.com/Ulqiora/Route256Project/internal/repository/pick_point/cli"
	"github.com/Ulqiora/Route256Project/internal/service/order/http/courier"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
	"github.com/Ulqiora/Route256Project/test/postgresql"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestReturnOrderToCourierIntegrate(t *testing.T) {
	var (
		err      error = nil
		ctx            = context.Background()
		changers       = map[model.TypePacking]order_changers.ChangerOrder{
			model.TypePackage: order_changers.ChangerOrderPackage{},
		}
	)
	configfile := http_service.LoadConfig("/home/ccnc/GolandProjects/homework/config/config-testing.example.yaml")
	database, err := postgresql.NewTestDatabase(ctx, configfile)
	require.Nil(t, err)
	defer func() {
		database.TruncateTable(ctx, "\"order\"")
		database.TruncateTable(ctx, "pickpoint")
	}()
	pickpoint_repo := pickpoint_storage.New(*database)
	id, err := pickpoint_repo.Create(ctx, pick_point_storage.PickPointDTO{
		Name:           "fds",
		Address:        "fds",
		ContactDetails: nil,
	})
	require.Nil(t, err)
	t.Run("ok", func(t *testing.T) {
		repository := order_storage.New(*database)
		initData := model.OrderInitData{
			CustomerID:  1,
			PickPointID: int64(id),
			ShelfLife:   jtime.TimeWrap(time.Now().Add(time.Hour)),
			Penny:       1000,
			Weight:      12,
			Type:        "Box",
		}
		idOrder, _ := repository.AcceptOrder(ctx, initData.MapToDTO())

		service := courier.New(repository, changers)
		router := mux.NewRouter()
		router.HandleFunc("/order/{id:[0-9]+}", service.ReturnOrderToCourier).Methods(http.MethodDelete)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/order/%d", idOrder), nil)
		router.ServeHTTP(rec, req)
		require.Equal(t, rec.Code, http.StatusOK)
	})
}
