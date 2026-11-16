package application

import "pos/domain"

type CalculateFinalGradeUseCase struct {
	studentsRepository domain.StudentRepository
	rubricRepository   domain.RubricRepository
	gradeCalculator    domain.GradeCalculator
}

func NewCalculateFinalGradeUseCase(studentsRepository domain.StudentRepository,
	rubricRepository domain.RubricRepository,
	gradeCalculator domain.GradeCalculator,
) CalculateFinalGradeUseCase {
	return CalculateFinalGradeUseCase{studentsRepository: studentsRepository, rubricRepository: rubricRepository, gradeCalculator: gradeCalculator}
}

func (uc CalculateFinalGradeUseCase) Execute(studentID string) (domain.FinalGrade, error) {
	student, err := uc.studentsRepository.FindByID(studentID)
	if err != nil {
		return domain.FinalGrade{}, err
	}
	rubric, err := uc.rubricRepository.Load()
	if err != nil {
		return domain.FinalGrade{}, err
	}
	return uc.gradeCalculator.Calculate(student, rubric)
}
