package domain

import (
	"fmt"
	"testing"
)

func TestWeight(t *testing.T) {
	type testCaseWeight struct {
		value float64
		isErr bool
	}

	testCasesWeight := []testCaseWeight{
		{0.3, false},
		{1, false},
		{0.0001, false},
		{0, true},
		{-0.1, true},
		{1.0001, true},
	}

	for i, weight := range testCasesWeight {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			_, err := NewWeight(weight.value)
			if (err != nil) != weight.isErr {
				t.Errorf("Weight(%f): ожидается ли ошибка: %t, получена ли ошибка: %t", weight.value, weight.isErr, (err != nil))
			}
		})
	}
}
