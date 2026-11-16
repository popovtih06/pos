package domain

import (
	"fmt"
	"math"
	"testing"
)

func TestScore(t *testing.T) {
	type testCaseNewScore struct {
		value int
		max   int
		isErr bool
	}

	type testCaseRatio struct {
		value int
		max   int
		want  float64
	}

	testCasesNewScore := []testCaseNewScore{
		{5, 10, false},
		{0, 10, false},
		{10, 10, false},
		{-1, 10, true},
		{11, 10, true},
		{5, 0, true},
		{5, -1, true},
	}

	testCasesRatio := []testCaseRatio{
		{16, 20, 0.8},
		{0, 1, 0},
		{3, 4, 0.75},
	}

	for i, score := range testCasesNewScore {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			_, err := NewScore(score.value, score.max)
			if (err != nil) != score.isErr {
				t.Errorf("NewScore(%d, %d): ожидается ли ошибка: %t, получена ли ошибка: %t", score.value, score.max, score.isErr, (err != nil))
			}
		})
	}

	for i, score := range testCasesRatio {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			sc, _ := NewScore(score.value, score.max)
			if math.Abs(sc.Ratio()-score.want) > 1e-9 {
				t.Errorf("Ratio от NewScore(%d, %d) ожидается %f, получено %f", score.value, score.max, score.want, sc.Ratio())
			}
		})
	}

}
