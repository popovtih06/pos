package infrastructure

import (
	"pos/domain"
)

type evaluationDTO struct {
	CriterionName string
	Value         int
	Max           int
}

type studentDTO struct {
	ID          string
	Name        string
	Evaluations []evaluationDTO
}

type criterionDTO struct {
	Name     string
	MaxScore int
	Weight   float64
}

type rubricDTO struct {
	Criteria []criterionDTO
}

func studentToDTO(student domain.Student) studentDTO {
	var stu studentDTO
	stu.ID = student.ID()
	stu.Name = student.Name()

	for k, v := range student.Evaluations() {
		var ev evaluationDTO
		ev.CriterionName = k
		ev.Value = v.Score().Value()
		ev.Max = v.Score().Max()
		stu.Evaluations = append(stu.Evaluations, ev)
	}
	return stu
}

func dtoToStudent(dto studentDTO, rubric domain.Rubric) (domain.Student, error) {
	stu, err := domain.NewStudent(dto.ID, dto.Name)
	if err != nil {
		return domain.Student{}, err
	}

	for _, ev := range dto.Evaluations {
		score, err := domain.NewScore(ev.Value, ev.Max)
		if err != nil {
			return domain.Student{}, err
		}
		evaluation, err := domain.NewEvaluation(rubric, ev.CriterionName, score)
		if err != nil {
			return domain.Student{}, err
		}
		if err := stu.SetEvaluation(evaluation, rubric); err != nil {
			return domain.Student{}, err
		}
	}
	return stu, nil
}

func rubricToDTO(rubric domain.Rubric) rubricDTO {
	var result rubricDTO
	var criteria []criterionDTO
	for k, v := range rubric.Criteria() {
		var cr criterionDTO
		cr.Name = k
		cr.MaxScore = v.MaxScore()
		cr.Weight = v.Weight().Value()
		criteria = append(criteria, cr)
	}
	result.Criteria = criteria
	return result
}

func dtoToRubric(dto rubricDTO) (domain.Rubric, error) {
	rubric := domain.NewRubric()
	for _, critDTO := range dto.Criteria {
		weight, err := domain.NewWeight(critDTO.Weight)
		if err != nil {
			return domain.Rubric{}, err
		}
		criterion, err := domain.NewCriterion(critDTO.Name, critDTO.MaxScore, weight)
		if err != nil {
			return domain.Rubric{}, err
		}
		err = rubric.AddCriterion(criterion)
		if err != nil {
			return domain.Rubric{}, err
		}
	}
	return rubric, nil
}
