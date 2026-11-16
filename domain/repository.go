package domain

import "errors"

var ErrStudentNotFound = errors.New("студент не найден")

type StudentRepository interface {
	Save(s Student) error
	FindByID(id string) (Student, error)
	All() ([]Student, error)
}

type RubricRepository interface {
	Save(r Rubric) error
	Load() (Rubric, error)
}
