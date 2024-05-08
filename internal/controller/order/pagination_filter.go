package order

import (
	"strconv"

	"homework/internal/controller"
)

type paginationFilter struct {
	Page  uint64
	Limit uint64
}

func (f *paginationFilter) loadFromRequest(values controller.ValuesView) error {
	var err error
	if values.Has("page") {
		f.Page, err = strconv.ParseUint(values.Get("page"), 10, 64)
		if err != nil {

			return err
		}
		f.Page--
	} else {
		f.Page = 0
	}
	if values.Has("limit") {
		f.Limit, err = strconv.ParseUint(values.Get("limit"), 10, 64)
		if err != nil {
			return err
		}
	} else {
		f.Limit = 10
	}
	return nil
}
