package infrastructure

import (
	"encoding/json"
	"errors"
	"os"
	"pos/domain"
)

type JSONStudentRepository struct {
	path             string
	rubricRepository domain.RubricRepository
}

func NewJSONStudentRepository(path string, rubricRepository domain.RubricRepository) JSONStudentRepository {
	return JSONStudentRepository{path: path, rubricRepository: rubricRepository}
}

func (rep JSONStudentRepository) All() ([]domain.Student, error) {
	dto, err := rep.readAll()
	if err != nil {
		return []domain.Student{}, err
	}
	rubric, err := rep.rubricRepository.Load()
	if err != nil {
		return nil, err
	}
	var dom []domain.Student
	for _, stuDto := range dto {
		stu, err := dtoToStudent(stuDto, rubric)
		if err != nil {
			return []domain.Student{}, err
		}
		dom = append(dom, stu)
	}
	return dom, nil
}

func (rep JSONStudentRepository) FindByID(id string) (domain.Student, error) {
	students, err := rep.All()
	if err != nil {
		return domain.Student{}, err
	}
	for _, student := range students {
		if student.ID() == id {
			return student, nil
		}
	}
	return domain.Student{}, domain.ErrStudentNotFound
}

func (rep JSONStudentRepository) Save(s domain.Student) error {
	rubric, err := rep.rubricRepository.Load()
	if err != nil {
		return err
	}
	if _, err := dtoToStudent(studentToDTO(s), rubric); err != nil {
		return err
	}
	dtoAll, err := rep.readAll()
	if err != nil {
		return err
	}
	studentDTO := studentToDTO(s)
	for i := range dtoAll {
		if dtoAll[i].ID == studentDTO.ID {
			dtoAll[i] = studentDTO
			return rep.writeAll(dtoAll)
		}
	}
	dtoAll = append(dtoAll, studentDTO)
	return rep.writeAll(dtoAll)
}

func (rep JSONStudentRepository) writeAll(list []studentDTO) error {
	raw, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rep.path, raw, 0644)
}

func (rep JSONStudentRepository) readAll() ([]studentDTO, error) {
	data, err := os.ReadFile(rep.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return []studentDTO{}, err
	}
	var list []studentDTO
	err = json.Unmarshal(data, &list)
	if err != nil {
		return []studentDTO{}, err
	}
	return list, nil
}

var _ domain.StudentRepository = JSONStudentRepository{}
