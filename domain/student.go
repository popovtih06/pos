package domain

import (
	"errors"
	"strings"
)

type Student struct {
	id          string
	name        string
	evaluations map[string]Evaluation
}

func NewStudent(id, name string) (Student, error) {
	if strings.TrimSpace(id) == "" {
		return Student{}, errors.New("id студента не может быть пустым")
	}
	if strings.TrimSpace(name) == "" {
		return Student{}, errors.New("имя студента не может быть пустым")
	}
	return Student{
		id:          id,
		name:        name,
		evaluations: make(map[string]Evaluation),
	}, nil
}

func (s Student) ID() string {
	return s.id
}
func (s Student) Name() string {
	return s.name
}

func (s *Student) Rename(newName string) error {
	if strings.TrimSpace(newName) == "" {
		return errors.New("имя студента не может быть пустым")
	}
	s.name = newName
	return nil
}

func (s *Student) SetEvaluation(e Evaluation, rubric Rubric) error {
	if _, err := NewEvaluation(rubric, e.CriterionName(), e.Score()); err != nil {
		return err
	}
	if s.evaluations == nil {
		s.evaluations = make(map[string]Evaluation)
	}
	s.evaluations[e.CriterionName()] = e
	return nil
}

func (s *Student) RemoveEvaluation(criterionName string) {
	delete(s.evaluations, criterionName)
}

func (s Student) Evaluations() map[string]Evaluation {
	out := make(map[string]Evaluation, len(s.evaluations))
	for k, v := range s.evaluations {
		out[k] = v
	}
	return out
}
