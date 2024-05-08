package model

import (
	"homework/internal/repository"
)

type ContactDetail struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func (c *ContactDetail) MapToDTO() repository.ContactDetailDTO {
	return repository.ContactDetailDTO{
		Type:   c.Type,
		Detail: c.Detail,
	}
}

func (c *ContactDetail) LoadFromDTO(dto repository.ContactDetailDTO) ContactDetail {
	*c = ContactDetail{
		Type:   dto.Type,
		Detail: dto.Detail,
	}
	return *c
}

func LoadContactDetailsFromDTO(dtos []repository.ContactDetailDTO) []ContactDetail {
	result := make([]ContactDetail, len(dtos))
	for i := range dtos {
		result[i].LoadFromDTO(dtos[i])
	}
	return result
}
