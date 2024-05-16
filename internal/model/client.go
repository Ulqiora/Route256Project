package model

import (
	"errors"

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
	err := orderdto.ID.Set(o.ID)
	return orderdto, err
}

func (o *Client) LoadFromDTO(dto *repository.ClientDTO) (Client, error) {
	if dto == nil {
		return Client{}, errors.New("dto is nil")
	}
	value, err := dto.ID.Value()
	if err != nil {
		return Client{}, err
	}
	*o = Client{
		ID:   value.(string),
		Name: dto.Name,
	}
	return *o, nil
}

func LoadClientsFromDTO(orders []repository.ClientDTO) ([]Client, error) {
	result := make([]Client, len(orders))
	for i := range orders {
		_, err := result[i].LoadFromDTO(&orders[i])
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

func (o *Client) LoadFromGrpcModel(dto *api.Client) (*Client, error) {
	if dto == nil {
		return nil, errors.New("client info is nil")
	}
	if dto.ID != nil {
		o.ID = dto.ID.Value
	}
	o.Name = dto.Name
	return o, nil
}
