package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/repository"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
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

func (o *OrderInitData) MapToDTO() (repository.OrderDTO, error) {
	orderdto := repository.OrderDTO{}
	err := orderdto.CustomerID.Set(o.CustomerID)
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
	return orderdto, nil
}

func (o *OrderInitData) LoadFromGrpcModel(dto *api.OrderInitData) error {
	if dto == nil {
		return errors.New("client info is nil")
	}
	*o = OrderInitData{
		CustomerID:  dto.Customer_ID.Value,
		PickPointID: dto.Pickpoint_ID.Value,
		ShelfLife:   jtime.TimeWrap(dto.ShelfTime.AsTime()),
		Penny:       dto.Penny,
		Weight:      dto.Weight,
		Type:        TypePacking(dto.TypePacking),
	}
	return nil
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
