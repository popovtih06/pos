package domain

type GradeCalculator interface {
	Calculate(s Student, r Rubric) (FinalGrade, error)
}

type WeightedCalculator struct{}

func (c WeightedCalculator) Calculate(s Student, r Rubric) (FinalGrade, error) {
	var sum float64
	for k, v := range s.Evaluations() {
		if _, err := NewEvaluation(r, k, v.Score()); err != nil {
			return FinalGrade{}, err
		}
		cr, err := r.FindCriterion(k)
		if err != nil {
			return FinalGrade{}, err
		}
		sum += v.Score().Ratio() * cr.Weight().Value()
	}
	return NewFinalGrade(sum * 10)
}
