//go:build go1.23

package code

import (
	"fmt"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTypeAlias(name string, t types.Type) *types.Alias {
	var nopos token.Pos
	return types.NewAlias(types.NewTypeName(nopos, nil, name, nil), t)
}

func TestUnalias(t *testing.T) {
	type aTest struct {
		input    types.Type
		expected string
	}

	intAlias := createTypeAlias("i", types.Typ[types.Int])
	intNestedAlias := createTypeAlias("ii", intAlias)
	pointerAlias := createTypeAlias("p", types.NewPointer(intNestedAlias))
	pointerNestedAlias := createTypeAlias("pp", pointerAlias)

	intPointer := types.NewPointer(types.Typ[types.Int])
	intPointerAlias := createTypeAlias("ip", intPointer)

	theTests := []aTest{
		{pointerNestedAlias, "*int"},
		{intNestedAlias, "int"},
		{intPointer, "*int"},
		{intPointerAlias, "*int"},
		{types.Typ[types.Int], "int"},
	}

	for _, at := range theTests {
		t.Run(fmt.Sprintf("unalias-%s", at.input), func(t *testing.T) {
			require.Equal(t, at.expected, Unalias(at.input).String())
		})
	}
}
