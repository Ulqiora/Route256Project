package repository

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type OrderDTO struct {
	ID          pgtype.UUID        `json:"id" db:"id"`                     // ID Заказа
	CustomerID  pgtype.UUID        `json:"id_customer" db:"id_customer"`   // ID Заказчика
	PickPointID pgtype.UUID        `json:"id_pickpoint" db:"id_pickpoint"` // ID Заказчика
	ShelfLife   pgtype.Timestamptz `json:"shelf_life" db:"shelf_life"`     // Срок хранения заказа на ПВЗ
	TimeCreated pgtype.Timestamptz `json:"time_created" db:"time_created"` // Дата получения заказа ПВЗ
	DateReceipt pgtype.Timestamptz `json:"date_receipt" db:"date_receipt"` // Дата получения заказа Клиентом

	Penny  pgtype.Numeric `json:"price" db:"price"`
	Weight pgtype.Numeric `json:"weight" db:"weight"`

	State string `json:"id_state" db:"id_state"` // Состояние заказа
}

func (o *OrderDTO) LoadFromRow(row pgx.Row) error {
	err := row.Scan(o)
	return err
}
