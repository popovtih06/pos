package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pos/application"
	"pos/domain"
	"pos/infrastructure"
	"strconv"
	"strings"

	"go.uber.org/dig"
)

type App struct {
	listGrades          application.ListGradesUseCase
	addStudent          application.AddStudentUseCase
	addCriterion        application.AddCriterionUseCase
	setGrade            application.SetGradeUseCase
	removeGrade         application.RemoveGradeUseCase
	listStudents        application.ListStudentsUseCase
	listCriteria        application.ListCriteriaUseCase
	calculateFinalGrade application.CalculateFinalGradeUseCase
	exportGrades        application.ExportGradesUseCase
}

func NewApp(listGrades application.ListGradesUseCase, addStudent application.AddStudentUseCase,
	addCriterion application.AddCriterionUseCase,
	setGrade application.SetGradeUseCase,
	removeGrade application.RemoveGradeUseCase,
	listStudents application.ListStudentsUseCase,
	listCriteria application.ListCriteriaUseCase,
	calculateFinalGrade application.CalculateFinalGradeUseCase,
	exportGrades application.ExportGradesUseCase,
) App {
	return App{
		listGrades:          listGrades,
		addStudent:          addStudent,
		addCriterion:        addCriterion,
		setGrade:            setGrade,
		removeGrade:         removeGrade,
		listStudents:        listStudents,
		listCriteria:        listCriteria,
		calculateFinalGrade: calculateFinalGrade,
		exportGrades:        exportGrades,
	}
}

func ask(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(scanner.Text())
}

func askInt(scanner *bufio.Scanner, prompt string, errStr string) (int, error) {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return 0, errors.New("не удалось прочитать")
	}
	num, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return 0, errors.New(errStr)
	}
	return num, nil
}

func askFloat(scanner *bufio.Scanner, prompt string, errStr string) (float64, error) {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return 0, errors.New("не удалось прочитать")
	}
	num, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
	if err != nil {
		return 0, errors.New(errStr)
	}
	return num, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	container := dig.New()
	must(container.Provide(func(rubricRepository domain.RubricRepository) domain.StudentRepository {
		return infrastructure.NewJSONStudentRepository("students.json", rubricRepository)
	}))
	must(container.Provide(func() domain.RubricRepository {
		return infrastructure.NewJSONRubricRepository("rubrics.json")
	}))
	must(container.Provide(func() domain.GradeCalculator {
		return domain.NewRoundingCalculator(domain.WeightedCalculator{})
	}))
	must(container.Provide(func() application.GradeExporter {
		return infrastructure.ExtensionExporter{}
	}))
	must(container.Provide(application.NewAddStudentUseCase))
	must(container.Provide(application.NewAddCriterionUseCase))
	must(container.Provide(application.NewSetGradeUseCase))
	must(container.Provide(application.NewRemoveGradeUseCase))
	must(container.Provide(application.NewListStudentsUseCase))
	must(container.Provide(application.NewListCriteriaUseCase))
	must(container.Provide(application.NewCalculateFinalGradeUseCase))
	must(container.Provide(application.NewExportGradesUseCase))
	must(container.Provide(application.NewListGradesUseCase))
	must(container.Provide(NewApp))
	must(container.Invoke(run))
}

func run(app App) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("0. Выход")
		fmt.Println("1. Список студентов")
		fmt.Println("2. Добавить студента")
		fmt.Println("3. Список критериев")
		fmt.Println("4. Добавить критерий")
		fmt.Println("5. Выставить оценку")
		fmt.Println("6. Отменить оценку")
		fmt.Println("7. Итог по студенту")
		fmt.Println("8. Выгрузить итоги")
		fmt.Println("9. Оценки студента")
		fmt.Print("Выберите пункт: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())
		switch choice {
		case "9":
			id := ask(scanner, "Введите id студента: ")
			grades, err := app.listGrades.Execute(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(grades) == 0 {
				fmt.Println("Оценок пока нет")
			}
			for _, grade := range grades {
				fmt.Printf("%s: %d/%d\n", grade.CriterionName(), grade.Score().Value(), grade.Score().Max())
			}

		case "8":
			path := ask(scanner, "Введите путь к файлу: ")
			err := app.exportGrades.Execute(path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Итоги выгружены")
		case "7":
			id := ask(scanner, "Введите id студента: ")
			finalGrade, err := app.calculateFinalGrade.Execute(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Итоговая оценка студента: %.1f\n", finalGrade.Value())
		case "6":
			id := ask(scanner, "Введите id студента: ")
			name := ask(scanner, "Введите название критерия: ")
			err := app.removeGrade.Execute(id, name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Оценка отменена")
		case "5":
			id := ask(scanner, "Введите id студента: ")
			name := ask(scanner, "Введите название критерия: ")
			value, err := askInt(scanner, "Введите балл: ", "Балл должен быть целым числом")
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = app.setGrade.Execute(id, name, value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Оценка выставлена")
		case "4":
			name := ask(scanner, "Введите название критерия: ")
			maxScore, err := askInt(scanner, "Введите максимальный балл: ", "Максимальный балл должен быть целым")
			if err != nil {
				fmt.Println(err)
				continue
			}
			weight, err := askFloat(scanner, "Введите вес: ", "Вес должен быть числом от 0 до 1")
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = app.addCriterion.Execute(name, maxScore, weight)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Критерий добавлен")
		case "3":
			rubric, err := app.listCriteria.Execute()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for name, criterion := range rubric.Criteria() {
				fmt.Println(name, criterion.MaxScore(), criterion.Weight().Value())
			}
		case "2":
			id := ask(scanner, "Введите id студента: ")
			name := ask(scanner, "Введите имя студента: ")
			err := app.addStudent.Execute(id, name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Студент добавлен")
		case "1":
			list, err := app.listStudents.Execute()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, student := range list {
				fmt.Println(student.ID(), student.Name())
			}
		case "0":
			return
		default:
			fmt.Println("Нет такого пункта")
		}
	}
}
