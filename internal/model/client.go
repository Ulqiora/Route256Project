package model

import (
	"errors"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/repository"
)

type Client struct {
	ID   string `json:"id"`   // ID Заказа
	Name string `json:"name"` // Состояние заказа
}

func (o *Client) MapToDTO() (repository.ClientDTO, error) {
	orderdto := repository.ClientDTO{
		Name: o.Name,
	}
	if o.ID != "" {
		err := orderdto.ID.Set(o.ID)
		if err != nil {
			return repository.ClientDTO{}, fmt.Errorf("failed to set client ID: %w", err)
		}
	}
	return orderdto, nil
}

func (o *Client) LoadFromDTO(dto *repository.ClientDTO) error {
	if dto == nil {
		return errors.New("dto is nil")
	}
	value, err := dto.ID.Value()
	if err != nil {
		return fmt.Errorf("failed to get client ID: %w", err)
	}
	*o = Client{
		ID:   value.(string),
		Name: dto.Name,
	}
	return nil
}

func LoadClientsFromDTO(orders []repository.ClientDTO) ([]Client, error) {
	result := make([]Client, len(orders))
	for i := range orders {
		err := result[i].LoadFromDTO(&orders[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (o *Client) MapToGrpcModel() *api.Client {
	client := &api.Client{
		ID: &api.UUID{
			Value: o.ID,
		},
		Name: o.Name,
	}
	return client
}

func (o *Client) LoadFromGrpcModel(dto *api.Client) error {
	if dto == nil {
		return errors.New("client info is nil")
	}
	if dto.ID != nil {
		o.ID = dto.ID.Value
	}
	o.Name = dto.Name
	return nil
}
