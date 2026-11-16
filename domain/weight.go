package domain

import (
	"errors"
	"math"
)

type Weight struct {
	value float64
}

func NewWeight(value float64) (Weight, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 || value > 1 {
		return Weight{}, errors.New("вес должен быть больше 0 и не больше 1")
	}
	return Weight{value: value}, nil
}

func (w Weight) Value() float64 {
	return w.value
}
