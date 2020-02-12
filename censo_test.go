package censo_test

import (
	"encoding/json"
	"testing"

	"github.com/shrotavre/censo"
	"github.com/stretchr/testify/assert"
)

type Dummy struct {
	FieldA string
	FieldB int
	FieldC DummyDeep
	FieldD int
}

type DummyDeep struct {
	DeepFieldD string
}

func TestCensorStruct(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		schemas := []censo.C{
			censo.CSim("FieldC/DeepFieldD", "X"),
			censo.CBas("FieldB"),
			censo.CSim("FieldA", "X"),
		}

		target := Dummy{
			FieldA: "real_value",
			FieldB: 1234,
			FieldC: DummyDeep{DeepFieldD: "real_value"},
			FieldD: 1234,
		}

		// act
		// st := time.Now()
		err := censo.Censor(&target, schemas)
		// fmt.Println("str", time.Since(st))

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target.FieldA)
		assert.Equal(t, 0, target.FieldB)
		assert.Equal(t, "X", target.FieldC.DeepFieldD)
		assert.Equal(t, 1234, target.FieldD)
	})
}

func TestCensorJSON(t *testing.T) {
	t.Run("ok_basic", func(t *testing.T) {
		// arrange
		jsontarget := `{"FieldA":"real_value","FieldB":1234,"FieldC":{"DeepFieldD":"real_value"},"FieldD":1234}`
		schemas := []censo.C{
			censo.CBas("FieldC/DeepFieldD"),
			censo.CBas("FieldB"),
			censo.CBas("FieldA"),
		}
		var jsonobjtarget map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &jsonobjtarget)
		assert.NoError(t, err)

		// act
		// st := time.Now()
		err = censo.Censor(&jsonobjtarget, schemas)
		child, childok := jsonobjtarget["FieldC"].(map[string]interface{})
		// fmt.Println("json", time.Since(st))

		// assert
		assert.Nil(t, err)
		assert.True(t, childok)
		assert.EqualValues(t, "", jsonobjtarget["FieldA"])
		assert.EqualValues(t, 0, jsonobjtarget["FieldB"])
		assert.EqualValues(t, 1234, jsonobjtarget["FieldD"])
		assert.EqualValues(t, "", child["DeepFieldD"])
	})

	t.Run("ok_simple", func(t *testing.T) {
		// arrange
		jsontarget := `{"FieldA":"real_value","FieldB":1234,"FieldC":{"DeepFieldD":"real_value"},"FieldD":1234}`
		schemas := []censo.C{
			censo.CSim("FieldC/DeepFieldD", "X"),
			censo.CSim("FieldB", 999),
			censo.CSim("FieldA", "X"),
		}
		var jsonobjtarget map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &jsonobjtarget)
		assert.NoError(t, err)

		// act
		// st := time.Now()
		err = censo.Censor(&jsonobjtarget, schemas)
		child, childok := jsonobjtarget["FieldC"].(map[string]interface{})
		// fmt.Println("json", time.Since(st))

		// assert
		assert.Nil(t, err)
		assert.True(t, childok)
		assert.EqualValues(t, "X", jsonobjtarget["FieldA"])
		assert.EqualValues(t, 999, jsonobjtarget["FieldB"])
		assert.EqualValues(t, 1234, jsonobjtarget["FieldD"])
		assert.EqualValues(t, "X", child["DeepFieldD"])
	})
}
