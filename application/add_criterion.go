package application

import "pos/domain"

type AddCriterionUseCase struct {
	rubricRepository domain.RubricRepository
}

func NewAddCriterionUseCase(rubricRepository domain.RubricRepository) AddCriterionUseCase {
	return AddCriterionUseCase{rubricRepository: rubricRepository}
}

func (uc AddCriterionUseCase) Execute(name string, maxScore int, weight float64) error {
	criterionWeight, err := domain.NewWeight(weight)
	if err != nil {
		return err
	}
	criterion, err := domain.NewCriterion(name, maxScore, criterionWeight)
	if err != nil {
		return err
	}
	rubric, err := uc.rubricRepository.Load()
	if err != nil {
		return err
	}
	err = rubric.AddCriterion(criterion)
	if err != nil {
		return err
	}
	return uc.rubricRepository.Save(rubric)
}
