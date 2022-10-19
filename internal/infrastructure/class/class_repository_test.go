package class_test

import (
	"context"
	"testing"

	"nory/internal/domain"
	. "nory/internal/infrastructure/class"

	"github.com/stretchr/testify/assert"
)

func TestClassRepository(t *testing.T) {
	t.Parallel()
	repos := []classRepository{
		{
			Name: "memory",
			R:    NewClassRepositoryMem(),
		},
	}

	for _, r := range repos {
		t.Run(r.Name, func(t *testing.T) {
			t.Parallel()
			t.Run("create", r.testCreate)
			t.Run("get by class id", r.testGet)
			t.Run("get by owner id", r.testGetByOwnerId)
			t.Run("update class", r.testUpdate)
			t.Run("delete", r.testDelete)
		})
	}
}

type classRepository struct {
	Name string
	R    domain.ClassRepository
}

func (cr classRepository) testCreate(t *testing.T) {
	testCases := []struct {
		Name  string
		Class domain.Class
		Err   error
	}{
		{"success", domain.Class{ClassId: "foo", OwnerId: "abel"}, nil},
		{"success", domain.Class{ClassId: "bar", OwnerId: "abel"}, nil},
		{"success", domain.Class{ClassId: "baz", OwnerId: "abelia"}, nil},
		{"duplicate", domain.Class{ClassId: "bar", OwnerId: "abel"}, domain.ErrClassExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := cr.R.CreateClass(context.Background(), &tc.Class)
			assert.Equal(t, tc.Err, err, "missmatch err")
		})
	}
}

func (cr classRepository) testDelete(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Err     error
	}{
		{"existing class", "foo", nil},
		{"existing class", "bar", nil},
		{"existing class", "baz", nil},
		{"unexisting class", "baz", nil},
		{"unexisting class", "baz", nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := cr.R.DeleteClass(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err, "missmatch err")
		})
	}
}

func (cr classRepository) testGet(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Err     error
	}{
		{"existing class", "foo", nil},
		{"existing class", "bar", nil},
		{"existing class", "baz", nil},
		{"unexisting class", "foo-bar", domain.ErrClassNotFound},
		{"unexisting class", "foo-baz", domain.ErrClassNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			c, err := cr.R.GetClass(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err)
			if tc.Err == nil {
				assert.Equal(t, tc.ClassId, c.ClassId, "unexpected class id")
			}
		})
	}
}

func (cr classRepository) testGetByOwnerId(t *testing.T) {
	testCases := []struct {
		Name    string
		OwnerId string
		Len     int
		Err     error
	}{
		{"exists", "abel", 2, nil},
		{"exists", "abelia", 1, nil},
		{"unexists", "abelia n", 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			classes, err := cr.R.GetByOwnerId(context.Background(), tc.OwnerId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.Equal(t, tc.Len, len(classes), "unexpected class received")
				for _, c := range classes {
					assert.Equal(t, tc.OwnerId, c.OwnerId, "unexpected class owner")
				}
			}
		})
	}
}

func (cr classRepository) testUpdate(t *testing.T) {
	testCases := []struct {
		Name  string
		Class domain.Class
		Err   error
	}{
		{"success", domain.Class{ClassId: "foo", Description: "foo"}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			prev, err := cr.R.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, nil, err, "unexpected error")

			err = cr.R.UpdateClass(context.Background(), &tc.Class)
			assert.Equal(t, tc.Err, err, "missmatch error")

			curr, err := cr.R.GetClass(context.Background(), tc.Class.ClassId)
			assert.Equal(t, nil, err, "unexpected error")

			if tc.Err == nil {
				assert.Equal(t, prev.ClassId, curr.ClassId, "should not update class id")
				assert.Equal(t, prev.OwnerId, curr.OwnerId, "should not update owner id")
				assert.Equal(t, tc.Class.Description, curr.Description, "should able update Description")
			}
		})
	}
}
