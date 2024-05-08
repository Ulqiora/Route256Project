package model

import (
	"encoding/json"
	"io"

	"homework/internal/repository"
)

type PickPoint struct {
	ID             int             `json:"id"`
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

func (p *PickPoint) MapToDTO() repository.PickPointDTO {
	return repository.PickPointDTO{
		ID:             p.ID,
		Name:           p.Name,
		Address:        p.Address,
		ContactDetails: nil,
	}
}

func (p *PickPoint) LoadFromDTO(dto repository.PickPointDTO) PickPoint {
	*p = PickPoint{
		ID:             dto.ID,
		Name:           dto.Name,
		Address:        dto.Address,
		ContactDetails: nil,
	}
	p.ContactDetails = make([]ContactDetail, len(dto.ContactDetails))
	for i, details := range dto.ContactDetails {
		p.ContactDetails[i].LoadFromDTO(details)
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
