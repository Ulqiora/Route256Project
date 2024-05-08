package order_changers

import (
	"homework/internal/model"
)

type ChangerOrderBox struct {
}

func (c ChangerOrderBox) Change(dto model.OrderInitData) (model.OrderInitData, error) {
	const minimalWeight = 3000

	if dto.Weight >= minimalWeight {
		return dto, ErrorHeavyWeight
	}
	dto.Penny += 2000
	return dto, nil
}
