package application

import (
	"pos/domain"
	"sort"
)

type ListGradesUseCase struct{ students domain.StudentRepository }

func NewListGradesUseCase(students domain.StudentRepository) ListGradesUseCase {
	return ListGradesUseCase{students: students}
}
func (uc ListGradesUseCase) Execute(id string) ([]domain.Evaluation, error) {
	student, err := uc.students.FindByID(id)
	if err != nil {
		return nil, err
	}
	grades := make([]domain.Evaluation, 0, len(student.Evaluations()))
	for _, grade := range student.Evaluations() {
		grades = append(grades, grade)
	}
	sort.Slice(grades, func(i, j int) bool { return grades[i].CriterionName() < grades[j].CriterionName() })
	return grades, nil
}
