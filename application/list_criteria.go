package application

import "pos/domain"

type ListCriteriaUseCase struct {
	rubricRepository domain.RubricRepository
}

func NewListCriteriaUseCase(rubricRepository domain.RubricRepository) ListCriteriaUseCase {
	return ListCriteriaUseCase{rubricRepository: rubricRepository}
}

func (uc ListCriteriaUseCase) Execute() (domain.Rubric, error) {
	return uc.rubricRepository.Load()
}
