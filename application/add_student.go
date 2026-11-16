package application

import (
	"errors"
	"pos/domain"
)

type AddStudentUseCase struct {
	studentsRepository domain.StudentRepository
}

func NewAddStudentUseCase(studentsRepository domain.StudentRepository) AddStudentUseCase {
	return AddStudentUseCase{studentsRepository: studentsRepository}
}

func (as AddStudentUseCase) Execute(studentID, studentName string) error {
	_, err := as.studentsRepository.FindByID(studentID)
	if err == nil {
		return errors.New("студент с таким id уже существует")
	}
	if !errors.Is(err, domain.ErrStudentNotFound) {
		return err
	}
	student, err := domain.NewStudent(studentID, studentName)
	if err != nil {
		return err
	}
	return as.studentsRepository.Save(student)
}
