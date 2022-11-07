package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"matheus.com/vgs/internal/model"
)

func TestEvaluatorEngine(t *testing.T) {
	engine := NewEvaluatorEngine()
	t.Run("cannot evaluate nil filter", func(t *testing.T) {
		var filter map[string]interface{}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("cannot evaluate filter with missing name", func(t *testing.T) {
		filter := map[string]interface{}{
			"city": "campinas",
		}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestFilterByCity(t *testing.T) {
	engine := NewEvaluatorEngine()
	t.Run(`filter byCity fails if "city" field is missing`, func(t *testing.T) {
		filter := map[string]interface{}{
			"filterName": "byCity",
		}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("filter byCity fails with invalid users type", func(t *testing.T) {
		filter := map[string]interface{}{
			"filterName": "byCity",
			"city":       []string{"campinas"},
			"Users":      3,
		}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("filter by city", func(t *testing.T) {
		filter := map[string]interface{}{
			"filterName": "byCity",
			"city":       []string{"campinas"},
			"Users": []model.User{
				{Name: "user1", City: "campinas"},
				{Name: "user2", City: "limeira"},
				{Name: "user3", City: "campinas"},
			},
		}
		res, err := engine.Evaluate(filter)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Len(t, res, 2)
	})
}

func TestFilterByLastLogin(t *testing.T) {
	engine := NewEvaluatorEngine()
	t.Run(`filter byLastLogin fails if "date" field is missing`, func(t *testing.T) {
		filter := map[string]interface{}{
			"filterName": "byLastLogin",
		}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("filter byLastLogin fails with invalid users type", func(t *testing.T) {
		now := time.Now().UTC()
		aMonthAgo := now.AddDate(0, -1, 0)
		filter := map[string]interface{}{
			"filterName": "byLastLogin",
			"date":       float64(aMonthAgo.Unix()),
			"Users":      3,
		}
		res, err := engine.Evaluate(filter)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("filter by last login date", func(t *testing.T) {
		now := time.Now().UTC()
		aMonthAgo := now.AddDate(0, -1, 0)
		twoMonthsAgo := now.AddDate(0, -2, 0)
		filter := map[string]interface{}{
			"filterName": "byLastLogin",
			"date":       float64(aMonthAgo.Unix()),
			"Users": []model.User{
				{Name: "user1", LastLogin: now},
				{Name: "user2", LastLogin: aMonthAgo},
				{Name: "user3", LastLogin: twoMonthsAgo},
			},
		}
		res, err := engine.Evaluate(filter)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Len(t, res, 1)
	})
}
