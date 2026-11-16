package domain

import (
	"errors"
	"math"
)

type FinalGrade struct {
	value float64
}

func NewFinalGrade(value float64) (FinalGrade, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 10 {
		return FinalGrade{}, errors.New("оценка должна быть от 0 до 10")
	}
	return FinalGrade{value: value}, nil
}

func (g FinalGrade) Value() float64 {
	return g.value
}
