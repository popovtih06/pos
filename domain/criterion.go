package domain

import (
	"errors"
	"strings"
)

type Criterion struct {
	name     string
	maxScore int
	weight   Weight
}

func NewCriterion(name string, maxScore int, weight Weight) (Criterion, error) {
	if strings.TrimSpace(name) == "" {
		return Criterion{}, errors.New("имя не может быть пустым")
	}
	if _, err := NewWeight(weight.Value()); err != nil {
		return Criterion{}, err
	}
	if maxScore <= 0 {
		return Criterion{}, errors.New("максимальный балл должен быть больше 0")
	}
	return Criterion{name: name, maxScore: maxScore, weight: weight}, nil
}
func (c Criterion) Name() string {
	return c.name
}
func (c Criterion) MaxScore() int {
	return c.maxScore
}
func (c Criterion) Weight() Weight {
	return c.weight
}
func (c *Criterion) ChangeWeight(w Weight) {
	c.weight = w
}
