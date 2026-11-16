package domain

import "math"

type RoundingCalculator struct {
	inner GradeCalculator
}

func NewRoundingCalculator(inner GradeCalculator) RoundingCalculator {
	return RoundingCalculator{inner: inner}
}

func (c RoundingCalculator) Calculate(s Student, r Rubric) (FinalGrade, error) {
	grade, err := c.inner.Calculate(s, r)
	if err != nil {
		return FinalGrade{}, err
	}
	return NewFinalGrade(math.Round(grade.Value()))
}
