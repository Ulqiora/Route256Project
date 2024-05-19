package model

import (
	"errors"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/repository"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
	"github.com/jackc/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Order struct {
	ID          string         `json:"id"`          // ID Заказа
	CustomerID  string         `json:"customer_id"` // ID Заказчика
	PickPointID string         `json:"pick_point_id"`
	ShelfLife   jtime.TimeWrap `json:"shelf_life"`   // Срок хранения заказа на ПВЗ
	TimeCreated jtime.TimeWrap `json:"time_created"` // Дата получения заказа ПВЗ
	DateReceipt jtime.TimeWrap `json:"date_receipt"` // Дата получения заказа Клиентом

	Penny  int64 `json:"price"`
	Weight int64 `json:"weight"`

	State string `json:"state"` // Состояние заказа
}

func (o *Order) MapToDTO() (repository.OrderDTO, error) {
	orderdto := repository.OrderDTO{
		State: o.State,
	}
	err := orderdto.ID.Set(o.ID)
	if err != nil {
		return repository.OrderDTO{}, fmt.Errorf("error set id for order: %v", err)
	}
	err = orderdto.CustomerID.Set(o.CustomerID)
	if err != nil {
		return repository.OrderDTO{}, fmt.Errorf("error set customer id for order: %v", err)
	}
	err = orderdto.PickPointID.Set(o.PickPointID)
	if err != nil {
		return repository.OrderDTO{}, fmt.Errorf("error set pickpoint id for order: %v", err)
	}
	err = orderdto.Penny.Set(o.Penny)
	if err != nil {
		return repository.OrderDTO{}, fmt.Errorf("error set penny for order: %v", err)
	}
	err = orderdto.Weight.Set(o.Weight)
	if err != nil {
		return repository.OrderDTO{}, fmt.Errorf("error set id for order: %v", err)
	}
	if o.ShelfLife.IsZero() {
		orderdto.ShelfLife.Status = pgtype.Null
	}
	if o.TimeCreated.IsZero() {
		orderdto.TimeCreated.Status = pgtype.Null
	}
	if o.DateReceipt.IsZero() {
		orderdto.DateReceipt.Status = pgtype.Null
	}
	return orderdto, nil
}

func (o *Order) LoadFromDTO(dto *repository.OrderDTO) error {
	if dto == nil {
		return errors.New("dto is nil")
	}
	*o = Order{
		Penny:  dto.Penny.Int.Int64(),
		Weight: dto.Weight.Int.Int64(),
		State:  dto.State,
	}
	value, err := dto.ID.Value()
	if err != nil {
		return fmt.Errorf("error get id for order: %v", err)
	}
	o.ID = value.(string)
	value, err = dto.CustomerID.Value()
	if err != nil {
		return fmt.Errorf("error get id for order: %v", err)
	}
	o.ID = value.(string)
	value, err = dto.PickPointID.Value()
	if err != nil {
		return fmt.Errorf("error get id for order: %v", err)
	}
	o.ID = value.(string)

	if dto.ShelfLife.Status == pgtype.Present {
		o.ShelfLife = jtime.TimeWrap(dto.ShelfLife.Time)
	}
	if dto.TimeCreated.Status == pgtype.Present {
		o.TimeCreated = jtime.TimeWrap(dto.TimeCreated.Time)
	}
	if dto.DateReceipt.Status == pgtype.Present {
		o.DateReceipt = jtime.TimeWrap(dto.DateReceipt.Time)
	}
	return nil
}

func LoadOrdersFromDTO(orders []repository.OrderDTO) ([]Order, error) {
	result := make([]Order, len(orders))
	for i := range orders {
		err := result[i].LoadFromDTO(&orders[i])
		return nil, err
	}
	return result, nil
}

func (o *Order) MapToGrpcModel() *api.Order {
	client := &api.Order{
		ID:           &api.UUID{Value: o.ID},
		Customer_ID:  &api.UUID{Value: o.CustomerID},
		Pickpoint_ID: &api.UUID{Value: o.PickPointID},
		ShelfTime:    timestamppb.New(*o.ShelfLife.Time()),
		TimeCreated:  timestamppb.New(*o.TimeCreated.Time()),
		DateReceipt:  timestamppb.New(*o.DateReceipt.Time()),
		Penny:        o.Penny,
		Weight:       o.Weight,
		State:        o.State,
	}
	return client
}

func (o *Order) LoadFromGrpcModel(dto *api.Order) error {
	if dto == nil {
		return errors.New("client info is nil")
	}
	*o = Order{
		ID:          dto.ID.Value,
		CustomerID:  dto.Customer_ID.Value,
		PickPointID: dto.Pickpoint_ID.Value,
		ShelfLife:   jtime.TimeWrap(dto.ShelfTime.AsTime()),
		TimeCreated: jtime.TimeWrap(dto.TimeCreated.AsTime()),
		DateReceipt: jtime.TimeWrap(dto.DateReceipt.AsTime()),
		Penny:       dto.Penny,
		Weight:      dto.Weight,
		State:       dto.State,
	}
	return nil
}

// Состояния заказа
const (
	EReadyToIssued = "ReadyToIssued" // Готов к получению
	EReturned      = "Returned"      // Состояние возвращенного заказа
	EReceived      = "Received"
	EDeleted       = "Deleted" // Состояние полученного заказа
)
