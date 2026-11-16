package domain

import (
	"math"
	"testing"
)

func TestNonFiniteValues(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewWeight(value); err == nil {
			t.Errorf("accepted weight %v", value)
		}
		if _, err := NewFinalGrade(value); err == nil {
			t.Errorf("accepted final grade %v", value)
		}
	}
}

func TestRubricInvariants(t *testing.T) {
	r := NewRubric()
	c := mustCriterion(t, "project", 20, .6)
	if err := r.AddCriterion(c); err != nil {
		t.Fatal(err)
	}
	if err := r.AddCriterion(c); err == nil {
		t.Fatal("duplicate accepted")
	}
	if err := r.AddCriterion(mustCriterion(t, "exam", 10, .5)); err == nil {
		t.Fatal("total weight above one accepted")
	}
	if err := r.AddCriterion(Criterion{}); err == nil {
		t.Fatal("zero criterion accepted")
	}
	copy := r.Criteria()
	delete(copy, "project")
	if _, err := r.FindCriterion("project"); err != nil {
		t.Fatal("aggregate changed through returned map")
	}
}

func TestEvaluationValidation(t *testing.T) {
	r := NewRubric()
	if err := r.AddCriterion(mustCriterion(t, "project", 20, 1)); err != nil {
		t.Fatal(err)
	}
	f := NewEvaluationFactory(r)
	for _, tc := range []struct {
		name  string
		value int
	}{{"missing", 1}, {"project", -1}, {"project", 21}} {
		if _, err := f.Create(tc.name, tc.value); err == nil {
			t.Fatalf("invalid evaluation accepted: %+v", tc)
		}
	}
	score, _ := NewScore(5, 10)
	if _, err := NewEvaluation(r, "project", score); err == nil {
		t.Fatal("mismatched maximum accepted")
	}
	s, _ := NewStudent("1", "Student")
	if err := s.SetEvaluation(Evaluation{}, r); err == nil {
		t.Fatal("zero evaluation accepted")
	}
	e, err := f.Create("project", 15)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetEvaluation(e, NewRubric()); err == nil {
		t.Fatal("foreign evaluation accepted")
	}
	if err := s.SetEvaluation(e, r); err != nil {
		t.Fatal(err)
	}
	copy := s.Evaluations()
	delete(copy, "project")
	if len(s.Evaluations()) != 1 {
		t.Fatal("student changed through returned map")
	}
	result, err := NewRoundingCalculator(WeightedCalculator{}).Calculate(s, r)
	if err != nil || result.Value() != 8 {
		t.Fatalf("rounding: %v, %v", result.Value(), err)
	}
	s.RemoveEvaluation("project")
	result, err = (WeightedCalculator{}).Calculate(s, r)
	if err != nil || result.Value() != 0 {
		t.Fatalf("removed grade still counted: %v, %v", result.Value(), err)
	}
}
