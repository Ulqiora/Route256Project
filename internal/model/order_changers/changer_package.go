package order_changers

import (
	"homework/internal/model"
)

type ChangerOrderPackage struct {
}

func (c ChangerOrderPackage) Change(dto model.OrderInitData) (model.OrderInitData, error) {
	const minimalWeight = 1000

	if dto.Weight >= minimalWeight {
		return dto, ErrorHeavyWeight
	}
	dto.Penny += 500
	return dto, nil
}
