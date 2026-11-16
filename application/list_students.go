package application

import "pos/domain"

type ListStudentsUseCase struct {
	studentsRepository domain.StudentRepository
}

func NewListStudentsUseCase(studentsRepository domain.StudentRepository) ListStudentsUseCase {
	return ListStudentsUseCase{studentsRepository: studentsRepository}
}

func (uc ListStudentsUseCase) Execute() ([]domain.Student, error) {
	return uc.studentsRepository.All()
}
