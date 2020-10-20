package censo_test

import (
	"encoding/json"
	"testing"

	"github.com/shrotavre/censo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Dummy struct {
	FieldA string
	FieldB string
	FieldC int
	FieldD int
	FieldE DummyDeep
}

type DummyDeep struct {
	DeepFieldD string
}

func TestCensor(t *testing.T) {
	schema := []censo.C{}

	t.Run("ok_universal_field_sign", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CSim("*", "X")}
		target := Dummy{
			FieldA: "real_value",
			FieldC: 1234,
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target.FieldA)
		assert.Equal(t, "X", target.FieldE.DeepFieldD)
		assert.Equal(t, 0, target.FieldC)
	})

	t.Run("err_type_map_invalid", func(t *testing.T) {
		// arrange
		var target map[int]interface{}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.NotNil(t, err)
		assert.Equal(t, censo.ErrNotCensorable, err)
	})
}

func TestPowerCensor(t *testing.T) {
	powerfunc := func(fieldname string, fieldval interface{}) (placeholder interface{}) {
		placeholder = fieldval

		if _, ok := placeholder.(string); ok {
			placeholder = fieldname
		} else if _, ok := placeholder.(int); ok {
			placeholder = 9999
		}

		return
	}

	t.Run("ok_type_struct", func(t *testing.T) {
		// arrange
		target := Dummy{
			FieldA: "real_value",
			FieldC: 1234,
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.PowerCensor(&target, powerfunc)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "FieldA", target.FieldA)
		assert.Equal(t, "FieldE/DeepFieldD", target.FieldE.DeepFieldD)
		assert.Equal(t, 9999, target.FieldC)
	})

	t.Run("ok_type_map", func(t *testing.T) {
		// arrange
		jsontarget := `{"FieldA":"real_value","FieldE":{"DeepFieldD":"real_deep_value"}}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.PowerCensor(&target, powerfunc)
		child, ok := target["FieldE"].(map[string]interface{})
		require.True(t, ok)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "FieldE/DeepFieldD", child["DeepFieldD"])
		assert.Equal(t, "FieldA", target["FieldA"])
	})
}
