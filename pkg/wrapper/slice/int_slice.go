package slice

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []uint64

// Set позволяет нам установить значение флага []int
func (i *IntSlice) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		intValue, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return err
		}
		*i = append(*i, intValue)
	}
	return nil
}

func (i *IntSlice) String() string {
	return fmt.Sprintf("%v", *i)
}
