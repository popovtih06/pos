package domain

import "errors"

type Evaluation struct {
	criterionName string
	score         Score
}

func NewEvaluation(rubric Rubric, criterionName string, score Score) (Evaluation, error) {
	if criterionName == "" {
		return Evaluation{}, errors.New("критерий не может быть пустым")
	}
	criterion, err := rubric.FindCriterion(criterionName)
	if err != nil {
		return Evaluation{}, err
	}
	if score.Max() != criterion.MaxScore() {
		return Evaluation{}, errors.New("максимум оценки не совпадает с критерием")
	}
	if _, err := NewScore(score.Value(), score.Max()); err != nil {
		return Evaluation{}, err
	}
	return Evaluation{criterionName: criterionName, score: score}, nil
}
func (e Evaluation) CriterionName() string {
	return e.criterionName
}
func (e Evaluation) Score() Score {
	return e.score
}
