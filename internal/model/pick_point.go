package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/repository"
)

type PickPoint struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Address        string          `json:"address"`
	ContactDetails []ContactDetail `json:"contact_details"`
}

func (p *PickPoint) Load(data io.Reader) error {
	body, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &p); err != nil {
		return err
	}
	return nil
}

func (p *PickPoint) MapToDTO() (repository.PickPointDTO, error) {
	orderdto := repository.PickPointDTO{
		Name:           p.Name,
		Address:        p.Address,
		ContactDetails: make([]repository.ContactDetailDTO, 0, len(p.ContactDetails)),
	}
	if p.ID != "" {
		err := orderdto.ID.Set(p.ID)
		if err != nil {
			return repository.PickPointDTO{}, fmt.Errorf("error set id for order: %v", err)
		}
	}
	for i := range p.ContactDetails {
		orderdto.ContactDetails = append(orderdto.ContactDetails, repository.ContactDetailDTO{
			Type:   p.ContactDetails[i].Type,
			Detail: p.ContactDetails[i].Detail,
		})
	}
	return orderdto, nil
}

func (p *PickPoint) LoadFromDTO(dto repository.PickPointDTO) PickPoint {
	*p = PickPoint{
		Name:           dto.Name,
		Address:        dto.Address,
		ContactDetails: make([]ContactDetail, 0, len(dto.ContactDetails)),
	}
	value, err := dto.ID.Value()
	if err != nil {
		return PickPoint{}
	}
	p.ID = value.(string)
	p.ContactDetails = make([]ContactDetail, len(dto.ContactDetails))
	for i := range dto.ContactDetails {
		p.ContactDetails[i].LoadFromDTO(dto.ContactDetails[i])
	}
	return *p
}
func LoadPickPointsFromDTO(orders []repository.PickPointDTO) []PickPoint {
	result := make([]PickPoint, len(orders))
	for i := range orders {
		result[i].LoadFromDTO(orders[i])
	}
	return result
}

func (p *PickPoint) MapToGrpcModel() *api.PickPoint {
	client := &api.PickPoint{
		ID:             &api.UUID{Value: p.ID},
		Address:        p.Address,
		Name:           p.Name,
		ContactDetails: make([]*api.ContactDetails, 0, len(p.ContactDetails)),
	}
	for i := range p.ContactDetails {
		client.ContactDetails = append(client.ContactDetails, &api.ContactDetails{
			Type:   p.ContactDetails[i].Type,
			Detail: p.ContactDetails[i].Detail,
		})
	}
	return client
}

func (p *PickPoint) LoadFromGrpcModel(dto *api.PickPoint) error {
	if dto == nil {
		return errors.New("client info is nil")
	}
	*p = PickPoint{
		Name:           dto.Name,
		Address:        dto.Address,
		ContactDetails: make([]ContactDetail, 0, len(dto.ContactDetails)),
	}
	if dto.ID == nil {
		p.ID = dto.ID.Value
	}
	for i := range dto.ContactDetails {
		p.ContactDetails = append(p.ContactDetails, ContactDetail{
			Type:   dto.ContactDetails[i].Type,
			Detail: dto.ContactDetails[i].Detail,
		})
	}
	return nil
}
