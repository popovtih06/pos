package domain

import "errors"

type Rubric struct {
	criteria map[string]Criterion
}

func NewRubric() Rubric {
	return Rubric{criteria: make(map[string]Criterion)}
}

func (r *Rubric) AddCriterion(c Criterion) error {
	if _, err := NewCriterion(c.Name(), c.MaxScore(), c.Weight()); err != nil {
		return err
	}
	if r.criteria == nil {
		r.criteria = make(map[string]Criterion)
	}
	_, ok := r.criteria[c.Name()]
	if ok {
		return errors.New("критерий с таким именем уже есть")
	}
	newTotal := r.TotalWeight()
	newTotal += c.Weight().Value()
	if newTotal > 1+1e-9 {
		return errors.New("сумма весов больше 1")
	}
	r.criteria[c.Name()] = c
	return nil
}

func (r Rubric) FindCriterion(name string) (Criterion, error) {
	cr, ok := r.criteria[name]
	if !ok {
		return Criterion{}, errors.New("критерий \"" + name + "\" не найден в рубрикаторе")
	}
	return cr, nil

}
func (r Rubric) Criteria() map[string]Criterion {
	out := make(map[string]Criterion, len(r.criteria))
	for k, v := range r.criteria {
		out[k] = v
	}
	return out
}

func (r Rubric) TotalWeight() float64 {
	sum := 0.0
	for _, v := range r.criteria {
		sum += v.Weight().Value()
	}
	return sum
}
