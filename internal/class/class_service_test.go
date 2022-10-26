package class_test

import (
	"context"
	"strings"
	"testing"

	"nory/domain"
	. "nory/internal/class"
	classtask "nory/internal/class_task"

	"github.com/google/uuid"
	"github.com/rs/xid"
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
	t.Run("create class tasks", cst.createClassTask)
	t.Run("create, access and delete class", cst.classCreate)
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
		OwnerId: uuid.NewString(),
		Name:    "foo",
	}

	res, err := cst.classService.CreateClass(context.Background(), class)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)

	classRes, err := cst.classService.GetClassInfo(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, classRes.Code, "failed to create class")
	assert.Equal(t, class.Name, classRes.Data.Name)
	assert.Equal(t, class.ClassId, classRes.Data.ClassId)

	_, err = cst.classService.ClassRepository.GetClass(context.Background(), class.ClassId)
	assert.Nil(t, err)

	err = cst.classService.AccessClass(context.Background(), &domain.User{}, class.ClassId)
	assert.NotNil(t, err)

	err = cst.classService.AccessClass(context.Background(), &domain.User{UserId: class.OwnerId}, class.ClassId)
	assert.Nil(t, err)

	// delete existing class
	_, err = cst.classService.DeleteClass(context.Background(), class.ClassId)
	assert.Nil(t, err)

	// delete unexisting class
	_, err = cst.classService.DeleteClass(context.Background(), class.ClassId)
	assert.Nil(t, err)

	_, err = cst.classService.ClassRepository.GetClass(context.Background(), class.ClassId)
	assert.NotNil(t, err)

	testCases := []struct {
		class domain.Class
	}{
		{domain.Class{Name: ""}},
		{domain.Class{Name: "abelia narindi agsya - abel"}},
		{domain.Class{Name: "abel", Description: strings.Repeat("abelia", 50)}},
	}

	for _, tc := range testCases {
		tc.class.OwnerId = uuid.NewString()
		_, err := cst.classService.CreateClass(context.Background(), &tc.class)
		assert.NotNilf(t, err, "unexpected at %#+v", tc.class)

		_, err = cst.classService.DeleteClass(context.Background(), class.ClassId)
		assert.Nil(t, err)
	}
}

func (cst *classServiceTest) createClassTask(t *testing.T) {
	t.Parallel()
	task := &domain.ClassTask{
		ClassId:  xid.New().String(),
		AuthorId: uuid.NewString(),
	}

	res, err := cst.classService.CreateClassTask(context.Background(), task)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)

	testCases := []struct{ task domain.ClassTask }{
		{domain.ClassTask{Name: "abelia narindi agsya - abel"}},
		{domain.ClassTask{Description: strings.Repeat("abelia narindi agsya", 520)}},
	}

	for _, tc := range testCases {
		_, err := cst.classService.CreateClassTask(context.Background(), &tc.task)
		assert.NotNilf(t, err, "unexpected at %#+v", tc.task)
	}
}
