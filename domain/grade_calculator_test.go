package domain

import (
	"math"
	"testing"
)

func mustCriterion(t *testing.T, name string, maxScore int, weight float64) Criterion {
	t.Helper()
	testWeight, err := NewWeight(weight)
	if err != nil {
		t.Fatalf("не удалось создать вес %f, %v", weight, err)
	}
	criterion, err := NewCriterion(name, maxScore, testWeight)
	if err != nil {
		t.Fatalf("не удалось создать критерий, %v", err)
	}
	return criterion
}

func TestGradeCalculator(t *testing.T) {
	student, err := NewStudent("s-002", "Vasya")
	if err != nil {
		t.Fatalf("не удалось создать студента, %v", err)
	}
	rubric := NewRubric()
	err = rubric.AddCriterion(mustCriterion(t, "Домашние задания", 10, 0.3))
	if err != nil {
		t.Fatalf("не удалось добавить критерий, %v", err)
	}
	err = rubric.AddCriterion(mustCriterion(t, "Проект", 20, 0.4))
	if err != nil {
		t.Fatalf("не удалось добавить критерий, %v", err)
	}
	err = rubric.AddCriterion(mustCriterion(t, "Контрольная", 10, 0.3))
	if err != nil {
		t.Fatalf("не удалось добавить критерий, %v", err)
	}
	factory := NewEvaluationFactory(rubric)
	ev1, err := factory.Create("Домашние задания", 9)
	if err != nil {
		t.Fatalf("не удалось создать оценку, %v", err)
	}
	ev2, err := factory.Create("Проект", 16)
	if err != nil {
		t.Fatalf("не удалось создать оценку, %v", err)
	}
	ev3, err := factory.Create("Контрольная", 7)
	if err != nil {
		t.Fatalf("не удалось создать оценку, %v", err)
	}
	if err := student.SetEvaluation(ev1, rubric); err != nil {
		t.Fatal(err)
	}
	if err := student.SetEvaluation(ev2, rubric); err != nil {
		t.Fatal(err)
	}
	if err := student.SetEvaluation(ev3, rubric); err != nil {
		t.Fatal(err)
	}
	result, err := WeightedCalculator{}.Calculate(student, rubric)
	if err != nil {
		t.Fatalf("не удалось получить финальную оценку, %v", err)
	}
	if math.Abs(result.Value()-8.0) > 1e-9 {
		t.Fatalf("оценка студента неправильная. ожидаемый результат: %f, полученный результат: %f", 8.0, result.Value())
	}
}
