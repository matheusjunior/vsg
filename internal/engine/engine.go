package engine

import (
	"errors"
	"time"

	"github.com/antonmedv/expr"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/model"
)

type filter func(values map[string]interface{}) ([]model.User, error)
type filterMap map[string]filter

type EvaluatorEngine struct {
	filters filterMap
}

func NewEvaluatorEngine() *EvaluatorEngine {
	ee := &EvaluatorEngine{}
	ee.buildFilterMap()
	return ee
}

func (ee *EvaluatorEngine) Evaluate(filter map[string]interface{}) ([]model.User, error) {
	filterNameKey, ok := filter["filterName"]
	if !ok {
		return nil, errors.New(`missing "filterName" field`)
	}
	filterName := filterNameKey.(string)
	users, err := ee.filters[filterName](filter)
	if err != nil {
		logger.Logger().Error("expression evaluation failed", err)
		return nil, err
	}
	return users, nil
}

func (ee *EvaluatorEngine) buildFilterMap() {
	ee.filters = make(filterMap)
	ee.filters["byCity"] = ee.byCity
	ee.filters["byLastLogin"] = ee.byLastLogin
}

func (ee *EvaluatorEngine) byCity(values map[string]interface{}) ([]model.User, error) {
	if _, ok := values["city"]; !ok {
		return nil, errors.New(`missing "city" field`)
	}
	code := `filter(Users, {.City in city})`
	res, err := expr.Eval(code, values)
	if err != nil {
		return nil, err
	}
	return ee.assertResult(res)
}

func (ee *EvaluatorEngine) byLastLogin(values map[string]interface{}) ([]model.User, error) {
	value, ok := values["date"]
	if !ok {
		return nil, errors.New(`missing "date" field`)
	}
	date, _ := value.(float64)
	values["date"] = time.Unix(int64(date), 0)
	code := `filter(Users, {.LastLogin.Before(date)})`
	res, err := expr.Eval(code, values)
	if err != nil {
		return nil, err
	}
	return ee.assertResult(res)
}

func (ee *EvaluatorEngine) assertResult(res interface{}) ([]model.User, error) {
	r, ok := res.([]interface{})
	if !ok {
		return nil, errors.New("could not cast")
	}
	users := make([]model.User, 0, len(r))
	for _, t := range r {
		users = append(users, t.(model.User))
	}
	return users, nil
}
