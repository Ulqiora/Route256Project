package model

import (
	"github.com/jackc/pgtype"
	"homework/internal/repository"
	jtime "homework/pkg/wrapper/jsontime"
)

type Order struct {
	ID          uint64         `json:"id"`          // ID Заказа
	CustomerID  int64          `json:"customer_id"` // ID Заказчика
	PickPointID int64          `json:"pick_point_id"`
	ShelfLife   jtime.TimeWrap `json:"shelf_life"`   // Срок хранения заказа на ПВЗ
	TimeCreated jtime.TimeWrap `json:"time_created"` // Дата получения заказа ПВЗ
	DateReceipt jtime.TimeWrap `json:"date_receipt"` // Дата получения заказа Клиентом

	Penny  int64 `json:"price"`
	Weight int64 `json:"weight"`

	State string `json:"state"` // Состояние заказа
}

func (o *Order) MapToDTO() repository.OrderDTO {
	orderdto := repository.OrderDTO{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		PickPointID: o.PickPointID,

		Penny:  o.Penny,
		Weight: o.Weight,

		State: o.State,
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
	return orderdto
}

func (o *Order) LoadFromDTO(dto repository.OrderDTO) Order {
	*o = Order{
		ID:          dto.ID,
		CustomerID:  dto.CustomerID,
		PickPointID: dto.PickPointID,

		Penny:  dto.Penny,
		Weight: dto.Weight,

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
