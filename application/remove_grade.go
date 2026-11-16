package application

import "pos/domain"

type RemoveGradeUseCase struct {
	studentsRepository domain.StudentRepository
}

func NewRemoveGradeUseCase(studentsRepository domain.StudentRepository) RemoveGradeUseCase {
	return RemoveGradeUseCase{studentsRepository: studentsRepository}
}

func (uc RemoveGradeUseCase) Execute(studentID string, criterionName string) error {
	student, err := uc.studentsRepository.FindByID(studentID)
	if err != nil {
		return err
	}
	student.RemoveEvaluation(criterionName)
	return uc.studentsRepository.Save(student)
}
