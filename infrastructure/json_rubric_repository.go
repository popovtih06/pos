package infrastructure

import (
	"encoding/json"
	"errors"
	"os"
	"pos/domain"
)

type JSONRubricRepository struct {
	path string
}

func NewJSONRubricRepository(path string) JSONRubricRepository {
	return JSONRubricRepository{path: path}
}

func (rep JSONRubricRepository) Save(r domain.Rubric) error {
	dto := rubricToDTO(r)
	raw, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rep.path, raw, 0644)
}

func (rep JSONRubricRepository) Load() (domain.Rubric, error) {
	data, err := os.ReadFile(rep.path)
	if errors.Is(err, os.ErrNotExist) {
		return domain.NewRubric(), nil
	}
	if err != nil {
		return domain.Rubric{}, err
	}
	var rubDTO rubricDTO
	err = json.Unmarshal(data, &rubDTO)
	if err != nil {
		return domain.Rubric{}, err
	}
	return dtoToRubric(rubDTO)
}

var _ domain.RubricRepository = JSONRubricRepository{}
