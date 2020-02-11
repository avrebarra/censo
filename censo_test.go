package censo_test

import (
	"testing"

	"github.com/shrotavre/censo"
	"github.com/stretchr/testify/assert"
)

type Dummy struct {
	FieldA string
	FieldB int
	FieldC DummyDeep
}

type DummyDeep struct {
	DeepFieldD string
}

func TestCensor(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		schemas := []censo.C{
			censo.CBas("FieldB"),
			censo.CBas("FieldC"),
			censo.CSim("FieldA", "X"),
		}

		target := Dummy{
			FieldA: "real_value",
			FieldB: 1234,
			FieldC: DummyDeep{DeepFieldD: "real_value"},
		}

		// act
		err := censo.Censor(&target, schemas)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, "X", target.FieldA)
		assert.Equal(t, 0, target.FieldB)
	})
}
