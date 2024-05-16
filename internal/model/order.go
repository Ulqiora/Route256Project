package model

import (
	"github.com/Ulqiora/Route256Project/internal/repository"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
	"github.com/jackc/pgtype"
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

func (o *Order) MapToDTO() repository.OrderDTO {
	orderdto := repository.OrderDTO{
		ID:          pgtype.UUID{},
		CustomerID:  pgtype.UUID{},
		PickPointID: pgtype.UUID{},
		Penny:       pgtype.Numeric{},
		Weight:      pgtype.Numeric{},
		State:       o.State,
	}
	_ = orderdto.ID.Set(o.ID)
	_ = orderdto.CustomerID.Set(o.CustomerID)
	_ = orderdto.PickPointID.Set(o.PickPointID)
	_ = orderdto.Penny.Set(o.Penny)
	_ = orderdto.Weight.Set(o.Weight)
	if o.ShelfLife.IsZero() {
		orderdto.ShelfLife.Status = pgtype.Null
	}
	if o.TimeCreated.IsZero() {
		orderdto.TimeCreated.Status = pgtype.Null
	}
	if o.DateReceipt.IsZero() {
		orderdto.DateReceipt.Status = pgtype.Null
	}
	return orderdto
}

func (o *Order) LoadFromDTO(dto repository.OrderDTO) Order {
	*o = Order{
		ID:          string(dto.ID.Bytes[:]),
		CustomerID:  string(dto.CustomerID.Bytes[:]),
		PickPointID: string(dto.PickPointID.Bytes[:]),

		Penny:  dto.Penny.Int.Int64(),
		Weight: dto.Weight.Int.Int64(),

		State: dto.State,
	}
	if dto.ShelfLife.Status == pgtype.Present {
		o.ShelfLife = jtime.TimeWrap(dto.ShelfLife.Time)
	}
	if dto.TimeCreated.Status == pgtype.Present {
		o.TimeCreated = jtime.TimeWrap(dto.TimeCreated.Time)
	}
	if dto.DateReceipt.Status == pgtype.Present {
		o.DateReceipt = jtime.TimeWrap(dto.DateReceipt.Time)
	}
	return *o
}

func LoadOrdersFromDTO(orders []repository.OrderDTO) []Order {
	result := make([]Order, len(orders))
	for i := range orders {
		result[i].LoadFromDTO(orders[i])
	}
	return result
}

// Состояния заказа
const (
	EReadyToIssued = "ReadyToIssued" // Готов к получению
	EReturned      = "Returned"      // Состояние возвращенного заказа
	EReceived      = "Received"
	EDeleted       = "Deleted" // Состояние полученного заказа
)
