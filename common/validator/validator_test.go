package validator_test

import (
	"testing"

	"nory/common/response"
	. "nory/common/validator"
)

func TestValidator(t *testing.T) {
	type foo struct {
		Name string `validate:"required"`
		Nick string `validate:"min=1,max=3" json:"n1ck"`
	}

	type nestedFoo struct {
		Bar string `validate:"max=1"`
		Foo foo
	}

	testCase := []struct {
		name       string
		data       any
		errorsPath []string
	}{
		{
			name:       "pass validation",
			data:       foo{"bar", "baz"},
			errorsPath: []string{},
		},
		{
			name:       "failed validation",
			data:       foo{"", "baz-baz"},
			errorsPath: []string{"name", "n1ck"},
		},
		{
			name:       "failed nested validation",
			data:       nestedFoo{"bar", foo{"", "bazz"}},
			errorsPath: []string{"foo.name", "foo.n1ck", "bar"},
		},
	}

	for _, tc := range testCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateStruct(tc.data, "")
			if err == nil && len(tc.errorsPath) == 0 {
				return
			}
			res, ok := err.(*response.ResponseError)
			if !ok {
				t.Errorf("unknown error: %v", err)
				return
			}
			if len(tc.errorsPath) != len(res.Errors) {
				t.Error("missmatch error count")
			}
			for _, p := range tc.errorsPath {
				if _, ok := res.Errors[p]; !ok {
					t.Errorf("missing expected error %s", p)
				}
			}
		})
	}
}
