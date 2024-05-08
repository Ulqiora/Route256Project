package model

import (
	"encoding/json"
	"io"

	"homework/internal/repository"
	jtime "homework/pkg/wrapper/jsontime"
)

type TypePacking string

const (
	TypePackage TypePacking = "Package"
	TypeBox     TypePacking = "Box"
	TypeTape    TypePacking = "Tape"
)

type OrderInitData struct {
	CustomerID  int64          `json:"customer_id"` // ID Заказчика
	PickPointID int64          `json:"pick_point_id"`
	ShelfLife   jtime.TimeWrap `json:"shelf_life"` // Срок хранения заказа на ПВЗ

	Penny  int64 `json:"price"`
	Weight int64 `json:"weight"`

	Type TypePacking `json:"type"`
}

func (o *OrderInitData) MapToDTO() repository.OrderDTO {
	dto := repository.OrderDTO{
		CustomerID:  o.CustomerID,
		PickPointID: o.PickPointID,
		Penny:       o.Penny,
		Weight:      o.Weight,
	}
	_ = dto.ShelfLife.Set(o.ShelfLife.Time())
	return dto
}

func (o *OrderInitData) LoadFromRequest(r io.Reader) error {
	byteAll, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(byteAll, o); err != nil {
		return err
	}
	return nil
}
