package application

import "pos/domain"

type SetGradeUseCase struct {
	studentRepository domain.StudentRepository
	rubricRepository  domain.RubricRepository
}

func NewSetGradeUseCase(studentRepository domain.StudentRepository, rubricRepository domain.RubricRepository) SetGradeUseCase {
	return SetGradeUseCase{studentRepository: studentRepository, rubricRepository: rubricRepository}
}

func (uc SetGradeUseCase) Execute(studentID, criterionName string, value int) error {
	student, err := uc.studentRepository.FindByID(studentID)
	if err != nil {
		return err
	}
	rubric, err := uc.rubricRepository.Load()
	if err != nil {
		return err
	}
	factory := domain.NewEvaluationFactory(rubric)
	evaluation, err := factory.Create(criterionName, value)
	if err != nil {
		return err
	}
	if err := student.SetEvaluation(evaluation, rubric); err != nil {
		return err
	}
	return uc.studentRepository.Save(student)
}
