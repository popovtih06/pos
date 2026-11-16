package domain

type EvaluationFactory struct {
	rubric Rubric
}

func NewEvaluationFactory(rubric Rubric) EvaluationFactory {
	return EvaluationFactory{rubric: rubric}
}

func (f EvaluationFactory) Create(criterionName string, value int) (Evaluation, error) {
	criterion, err := f.rubric.FindCriterion(criterionName)
	if err != nil {
		return Evaluation{}, err
	}
	score, err := NewScore(value, criterion.MaxScore())
	if err != nil {
		return Evaluation{}, err
	}
	return NewEvaluation(f.rubric, criterionName, score)
}
