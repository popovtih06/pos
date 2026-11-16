package application_test

import (
	"errors"
	"os"
	"path/filepath"
	"pos/application"
	"pos/domain"
	"pos/infrastructure"
	"testing"
)

func TestGradingWorkflow(t *testing.T) {
	dir := t.TempDir()
	rubricRepo := infrastructure.NewJSONRubricRepository(filepath.Join(dir, "rubrics.json"))
	studentPath := filepath.Join(dir, "students.json")
	repo := infrastructure.NewJSONStudentRepository(studentPath, rubricRepo)
	add := application.NewAddStudentUseCase(repo)
	if err := add.Execute("1", "Иванов Иван"); err != nil {
		t.Fatal(err)
	}
	if err := add.Execute("1", "Другой студент"); err == nil {
		t.Fatal("duplicate student accepted")
	}
	if err := application.NewAddCriterionUseCase(rubricRepo).Execute("Проект", 20, 1); err != nil {
		t.Fatal(err)
	}
	set := application.NewSetGradeUseCase(repo, rubricRepo)
	if err := set.Execute("1", "Проект", 21); err == nil {
		t.Fatal("invalid score accepted")
	}
	if err := set.Execute("1", "Неизвестный", 1); err == nil {
		t.Fatal("unknown criterion accepted")
	}
	if err := set.Execute("1", "Проект", 15); err != nil {
		t.Fatal(err)
	}
	reloaded := infrastructure.NewJSONStudentRepository(studentPath, rubricRepo)
	grades, err := application.NewListGradesUseCase(reloaded).Execute("1")
	if err != nil || len(grades) != 1 || grades[0].Score().Value() != 15 {
		t.Fatalf("reload: %v, %v", grades, err)
	}
	calc := application.NewCalculateFinalGradeUseCase(reloaded, rubricRepo, domain.NewRoundingCalculator(domain.WeightedCalculator{}))
	result, err := calc.Execute("1")
	if err != nil || result.Value() != 8 {
		t.Fatalf("initial total: %v, %v", result.Value(), err)
	}
	if err := set.Execute("1", "Проект", 20); err != nil {
		t.Fatal(err)
	}
	result, err = calc.Execute("1")
	if err != nil || result.Value() != 10 {
		t.Fatalf("corrected total: %v, %v", result.Value(), err)
	}
	export := application.NewExportGradesUseCase(repo, rubricRepo, domain.WeightedCalculator{}, infrastructure.ExtensionExporter{})
	if err := export.Execute(filepath.Join(dir, "grades.csv")); err != nil {
		t.Fatal(err)
	}
	if err := export.Execute(filepath.Join(dir, "grades.json")); err != nil {
		t.Fatal(err)
	}
	if err := application.NewRemoveGradeUseCase(repo).Execute("1", "Проект"); err != nil {
		t.Fatal(err)
	}
	result, err = calc.Execute("1")
	if err != nil || result.Value() != 0 {
		t.Fatalf("removed total: %v, %v", result.Value(), err)
	}
	if _, err := repo.FindByID("missing"); !errors.Is(err, domain.ErrStudentNotFound) {
		t.Fatalf("not found: %v", err)
	}
}

func TestCorruptStorageIsNotOverwritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "students.json")
	original := []byte("broken json")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	repo := infrastructure.NewJSONStudentRepository(path, infrastructure.NewJSONRubricRepository(filepath.Join(dir, "rubrics.json")))
	if err := application.NewAddStudentUseCase(repo).Execute("1", "Student"); err == nil {
		t.Fatal("storage error ignored")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != string(original) {
		t.Fatalf("storage overwritten: %s, %v", data, err)
	}
}

func TestStorageRejectsInvalidEvaluation(t *testing.T) {
	for _, evaluation := range []string{
		`{"CriterionName":"missing","Value":5,"Max":10}`,
		`{"CriterionName":"project","Value":5,"Max":20}`,
		`{"CriterionName":"project","Value":11,"Max":10}`,
	} {
		dir := t.TempDir()
		rubricRepo := infrastructure.NewJSONRubricRepository(filepath.Join(dir, "rubrics.json"))
		if err := application.NewAddCriterionUseCase(rubricRepo).Execute("project", 10, 1); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "students.json")
		if err := os.WriteFile(path, []byte(`[{"ID":"1","Name":"Student","Evaluations":[`+evaluation+`]}]`), 0600); err != nil {
			t.Fatal(err)
		}
		repo := infrastructure.NewJSONStudentRepository(path, rubricRepo)
		if _, err := repo.All(); err == nil {
			t.Fatalf("invalid persisted evaluation accepted: %s", evaluation)
		}
	}
}
