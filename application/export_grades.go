package application

import "pos/domain"

type GradeRow struct {
	StudentName string
	Grade       float64
}

type GradeExporter interface {
	Export(rows []GradeRow, path string) error
}

type ExportGradesUseCase struct {
	studentsRepository domain.StudentRepository
	rubricRepository   domain.RubricRepository
	gradeCalculator    domain.GradeCalculator
	exporter           GradeExporter
}

func NewExportGradesUseCase(studentsRepository domain.StudentRepository,
	rubricRepository domain.RubricRepository,
	gradeCalculator domain.GradeCalculator,
	exporter GradeExporter,
) ExportGradesUseCase {
	return ExportGradesUseCase{
		studentsRepository: studentsRepository,
		rubricRepository:   rubricRepository,
		gradeCalculator:    gradeCalculator,
		exporter:           exporter,
	}
}

func (eg ExportGradesUseCase) Execute(path string) error {
	rub, err := eg.rubricRepository.Load()
	if err != nil {
		return err
	}
	allStu, err := eg.studentsRepository.All()
	if err != nil {
		return err
	}
	var res []GradeRow
	for _, stu := range allStu {
		var gr GradeRow
		gr.StudentName = stu.Name()
		fg, err := eg.gradeCalculator.Calculate(stu, rub)
		if err != nil {
			return err
		}
		gr.Grade = fg.Value()
		res = append(res, gr)
	}
	return eg.exporter.Export(res, path)
}
