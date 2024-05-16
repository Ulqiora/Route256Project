package model

import (
	"encoding/json"
	"io"

	"github.com/Ulqiora/Route256Project/internal/repository"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
	"github.com/jackc/pgtype"
)

type TypePacking string

const (
	TypePackage TypePacking = "Package"
	TypeBox     TypePacking = "Box"
	TypeTape    TypePacking = "Tape"
)

type OrderInitData struct {
	CustomerID  string         `json:"customer_id"` // ID Заказчика
	PickPointID string         `json:"pick_point_id"`
	ShelfLife   jtime.TimeWrap `json:"shelf_life"` // Срок хранения заказа на ПВЗ

	Penny  int64 `json:"price"`
	Weight int64 `json:"weight"`

	Type TypePacking `json:"type"`
}

func (o *OrderInitData) MapToDTO() repository.OrderDTO {
	dto := repository.OrderDTO{
		CustomerID:  pgtype.UUID{},
		PickPointID: pgtype.UUID{},
		Penny:       pgtype.Numeric{},
		Weight:      pgtype.Numeric{},
	}
	_ = dto.CustomerID.Set(o.CustomerID)
	_ = dto.PickPointID.Set(o.PickPointID)
	_ = dto.Weight.Set(o.Weight)
	_ = dto.Penny.Set(o.Penny)

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
