package class_test

import (
	"context"
	"testing"

	"nory/domain"
	. "nory/internal/class"
	classtask "nory/internal/class_task"

	"github.com/stretchr/testify/assert"
)

func TestClassService(t *testing.T) {
	t.Parallel()
	classService := ClassService{
		ClassRepository:     NewClassRepositoryMem(),
		ClassTaskRepository: classtask.NewClassTaskRepositoryMem(),
	}

	cst := classServiceTest{classService}

	t.Run("get class info", cst.classInfo)
	t.Run("get class tasks", cst.classTasks)
	t.Run("create class", cst.classCreate)
}

type classServiceTest struct {
	classService ClassService
}

func (cst classServiceTest) classInfo(t *testing.T) {
	t.Parallel()
	createClass := func() (c *domain.Class) {
		c = &domain.Class{}
		err := cst.classService.ClassRepository.CreateClass(context.Background(), c)
		assert.Nil(t, err)
		t.Cleanup(func() {
			cst.classService.ClassRepository.DeleteClass(context.Background(), c.ClassId)
		})
		return
	}

	_, err := cst.classService.GetClassInfo(context.Background(), "foobarbazqux")
	assert.ErrorContains(t, err, "can not find class with id \"foobarbazqux\"")
	classA := createClass()
	res, err := cst.classService.GetClassInfo(context.Background(), classA.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, classA, res.Data)
}

func (cst classServiceTest) classTasks(t *testing.T) {
	t.Parallel()
}

func (cst classServiceTest) classCreate(t *testing.T) {
	t.Parallel()

	class := &domain.Class{
		Name: "foo",
	}

	res, err := cst.classService.CreateClass(context.Background(), class)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)


	classRes, err := cst.classService.GetClassInfo(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, classRes.Code, "failed to create class")
	assert.Equal(t, class.Name, classRes.Data.Name)
	assert.Equal(t, class.ClassId, classRes.Data.ClassId)
}
