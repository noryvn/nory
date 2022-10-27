package validator_test

import (
	"strings"
	"testing"

	. "nory/common/validator"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	t.Parallel()
	type foo struct {
		Name string `validate:"required"`
		Nick string `validate:"min=1,max=3" json:"n1ck"`
	}

	type nestedFoo struct {
		Bar string `json:"-" validate:"max=1"`
		Foo foo
	}

	type bar struct {
		Userame string `validate:"username"`
	}

	testCase := []struct {
		name string
		data any
		err  bool
	}{
		{
			name: "pass validation",
			data: foo{"bar", "baz"},
			err:  false,
		},
		{
			name: "failed validation",
			data: foo{"", "baz-baz"},
			err:  true,
		},
		{
			name: "failed nested validation",
			data: nestedFoo{"bar", foo{"", "bazz"}},
			err:  true,
		},
		{
			name: "custom validator 'username'",
			data: bar{"abelia_narindi.agsya"},
			err:  false,
		},
		{
			name: "fail on forbidden character, validator 'username'",
			data: bar{"abel?"},
			err:  true,
		},
		{
			name: "fail on length, validator 'username'",
			data: bar{strings.Repeat("u", 21)},
			err:  true,
		},
	}

	for _, tc := range testCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			err := ValidateStruct(tc.data)
			if tc.err {
				assert.NotNil(t, err, "validate struct not return erro while error expected")
			} else {
				assert.Nil(t, err, "validate struct return error while not expected")
			}
		})
	}
}
