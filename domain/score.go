package domain

import "errors"

type Score struct {
	value int
	max   int
}

func NewScore(value, max int) (Score, error) {
	if value < 0 {
		return Score{}, errors.New("балл не может быть меньше 0")
	}
	if max <= 0 {
		return Score{}, errors.New("максимальное значение не может быть меньше или равно 0")
	}
	if value > max {
		return Score{}, errors.New("балл не может быть больше максимального значения")
	}
	return Score{value: value, max: max}, nil
}
func (s Score) Value() int {
	return s.value
}
func (s Score) Max() int {
	return s.max
}
func (s Score) Ratio() float64 {
	return float64(s.value) / float64(s.max)
}
