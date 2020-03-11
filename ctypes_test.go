package censo_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/shrotavre/censo"
)

func TestCensorCBas(t *testing.T) {
	schema := []censo.C{censo.CBas("FieldA")}

	t.Run("ok", func(t *testing.T) {
		// arrange
		target := Dummy{
			FieldA: "real_value",
			FieldB: "real_value",
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "", target.FieldA)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_map", func(t *testing.T) {
		// arrange
		jsontarget := `{"FieldA":"real_value","FieldE":"real_value"}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "", target["FieldA"])
		assert.Equal(t, "real_value", target["FieldE"])
	})

	t.Run("ok_type_nested", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CBas("FieldE/DeepFieldD")}
		target := Dummy{
			FieldB: "real_value",
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "", target.FieldE.DeepFieldD)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_nested_map", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CBas("FieldE/DeepFieldD")}
		jsontarget := `{"FieldA":"real_value","FieldE":{"DeepFieldD":"real_deep_value"}}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)
		child, ok := target["FieldE"].(map[string]interface{})
		require.True(t, ok)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "", child["DeepFieldD"])
		assert.Equal(t, "real_value", target["FieldA"])
	})

	t.Run("err_type_invalid", func(t *testing.T) {
		// arrange
		target := 1

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.NotNil(t, err)
	})
}

func TestCensorCSim(t *testing.T) {
	schema := []censo.C{censo.CSim("FieldA", "X")}

	t.Run("ok", func(t *testing.T) {
		// arrange
		target := Dummy{
			FieldA: "real_value",
			FieldB: "real_value",
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target.FieldA)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_map", func(t *testing.T) {
		jsontarget := `{"FieldA":"real_value","FieldE":"real_value"}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target["FieldA"])
		assert.Equal(t, "real_value", target["FieldE"])
	})

	t.Run("ok_type_nested", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CSim("FieldE/DeepFieldD", "X")}
		target := Dummy{
			FieldB: "real_value",
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target.FieldE.DeepFieldD)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_nested_map", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CSim("FieldE/DeepFieldD", "X")}
		jsontarget := `{"FieldA":"real_value","FieldE":{"DeepFieldD":"real_deep_value"}}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)
		child, ok := target["FieldE"].(map[string]interface{})
		require.True(t, ok)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", child["DeepFieldD"])
		assert.Equal(t, "real_value", target["FieldA"])
	})

	t.Run("ok_mismatch_type_should_empty", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CSim("FieldA", 1)}
		target := Dummy{
			FieldA: "real_value",
			FieldB: "real_value",
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "", target.FieldA)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("err_type_invalid", func(t *testing.T) {
		// arrange
		target := 1

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.NotNil(t, err)
	})
}

func TestCensorCFunc(t *testing.T) {
	replacer := func(i interface{}) (o interface{}) {
		o = i

		if v, ok := i.(string); ok && strings.Contains(v, "real") {
			o = strings.ReplaceAll(v, "real", "fake")
		}

		return
	}

	schema := []censo.C{
		censo.CFunc("FieldA", replacer),
	}

	t.Run("ok", func(t *testing.T) {
		// arrange
		target := Dummy{
			FieldA: "real_value",
			FieldB: "real_value",
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "fake_value", target.FieldA)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_map", func(t *testing.T) {
		// arrange
		jsontarget := `{"FieldA":"real_value","FieldE":"real_value"}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "fake_value", target["FieldA"])
		assert.Equal(t, "real_value", target["FieldE"])
	})

	t.Run("ok_type_nested", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CFunc("FieldE/DeepFieldD", replacer)}
		target := Dummy{
			FieldB: "real_value",
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "fake_deep_value", target.FieldE.DeepFieldD)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("ok_type_nested_map", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.CFunc("FieldE/DeepFieldD", replacer)}
		jsontarget := `{"FieldA":"real_value","FieldE":{"DeepFieldD":"real_deep_value"}}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)
		child, ok := target["FieldE"].(map[string]interface{})
		require.True(t, ok)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "fake_deep_value", child["DeepFieldD"])
		assert.Equal(t, "real_value", target["FieldA"])
	})

	t.Run("err_functype_invalid", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.C{"FieldE/DeepFieldD", func() {}}}
		target := Dummy{
			FieldB: "real_value",
			FieldE: DummyDeep{
				DeepFieldD: "real_deep_value",
			},
		}

		// act
		err := censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "real_deep_value", target.FieldE.DeepFieldD)
		assert.Equal(t, "real_value", target.FieldB)
	})

	t.Run("err_functype_invalid_map", func(t *testing.T) {
		// arrange
		schema := []censo.C{censo.C{"FieldE/DeepFieldD", func() {}}}
		jsontarget := `{"FieldA":"real_value","FieldE":{"DeepFieldD":"real_deep_value"}}`

		var target map[string]interface{}
		err := json.Unmarshal([]byte(jsontarget), &target)
		require.NoError(t, err)

		// act
		err = censo.Censor(&target, schema)

		// assert
		assert.Nil(t, err)
		child, ok := target["FieldE"].(map[string]interface{})
		require.True(t, ok)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "real_deep_value", child["DeepFieldD"])
		assert.Equal(t, "real_value", target["FieldA"])
	})
}
