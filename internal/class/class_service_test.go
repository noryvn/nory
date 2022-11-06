package class_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"nory/domain"
	. "nory/internal/class"
	classmember "nory/internal/class_member"
	classtask "nory/internal/class_task"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestClassService(t *testing.T) {
	t.Parallel()
	classService := ClassService{
		UserRepository:        user.NewUserRepositoryMem(),
		ClassRepository:       NewClassRepositoryMem(),
		ClassTaskRepository:   classtask.NewClassTaskRepositoryMem(),
		ClassMemberRepository: classmember.NewClassMemberRepositoryMem(),
	}

	cst := classServiceTest{classService}

	t.Run("get class info", cst.classInfo)
	t.Run("get class tasks", cst.classTasks)
	t.Run("create class tasks", cst.createClassTask)
	t.Run("create, access and delete class", cst.classCreate)
	t.Run("list member", cst.listMemberTask)
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

	class := &domain.Class{
		OwnerId: uuid.NewString(),
		Name:    "foo",
	}

	r, err := cst.classService.CreateClass(context.Background(), class)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.Code)

	tommorrow := time.Now().UTC().Add(24 * time.Hour)
	yesterday := tommorrow.Add(-2 * 24 * time.Hour)

	for i := 0; i < 5; i++ {
		err := cst.classService.ClassTaskRepository.CreateTask(context.Background(), &domain.ClassTask{
			ClassId:  class.ClassId,
			AuthorId: class.OwnerId,
			DueDate:  tommorrow,
		})
		assert.Nil(t, err)
		if i > 2 {
			err := cst.classService.ClassTaskRepository.CreateTask(context.Background(), &domain.ClassTask{
				ClassId:  class.ClassId,
				AuthorId: class.OwnerId,
				DueDate:  yesterday,
			})
			assert.Nil(t, err)
		}
	}

	for _, tc := range []struct {
		From time.Time
		To   time.Time
		Len  int
	}{
		{time.Time{}, time.Time{}, 5},
		{tommorrow, time.Time{}, 5},
		{yesterday, time.Time{}, 7},
		{yesterday, tommorrow, 2},
	} {
		res, err := cst.classService.GetClassTasks(context.Background(), class.ClassId, tc.From, tc.To)
		assert.Nil(t, err)
		assert.Equal(t, tc.Len, len(res.Data))
	}
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

	member := uuid.NewString()
	err = cst.classService.ClassMemberRepository.CreateMember(context.Background(), &domain.ClassMember{
		Level:   "member",
		UserId:  member,
		ClassId: class.ClassId,
	})
	assert.Nil(t, err)

	admin := uuid.NewString()
	err = cst.classService.ClassMemberRepository.CreateMember(context.Background(), &domain.ClassMember{
		Level:   "admin",
		UserId:  admin,
		ClassId: class.ClassId,
	})
	assert.Nil(t, err)

	classRes, err := cst.classService.GetClassInfo(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, classRes.Code, "failed to create class")
	assert.Equal(t, class.Name, classRes.Data.Name)
	assert.Equal(t, class.ClassId, classRes.Data.ClassId)

	_, err = cst.classService.ClassRepository.GetClass(context.Background(), class.ClassId)
	assert.Nil(t, err)

	err = cst.classService.AccessClass(context.Background(), uuid.NewString(), class.ClassId)
	assert.NotNil(t, err)

	err = cst.classService.AccessClass(context.Background(), class.OwnerId, class.ClassId)
	assert.Nil(t, err)

	err = cst.classService.AccessClass(context.Background(), member, class.ClassId)
	assert.NotNil(t, err)

	err = cst.classService.AccessClass(context.Background(), admin, class.ClassId)
	assert.Nil(t, err)

	// delete existing class
	_, err = cst.classService.DeleteClass(context.Background(), class.OwnerId, class.ClassId)
	assert.Nil(t, err)

	// delete unexisting class
	_, err = cst.classService.DeleteClass(context.Background(), class.OwnerId, xid.New().String())
	assert.NotNil(t, err)

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

		_, err = cst.classService.DeleteClass(context.Background(), class.OwnerId, class.ClassId)
		assert.NotNil(t, err)
	}
}

func (cst *classServiceTest) createClassTask(t *testing.T) {
	t.Parallel()
	class := &domain.Class{
		OwnerId: uuid.NewString(),
	}
	err := cst.classService.ClassRepository.CreateClass(context.Background(), class)
	assert.Nil(t, err)

	task := &domain.ClassTask{
		ClassId:  class.ClassId,
		AuthorId: uuid.NewString(),
	}

	res, err := cst.classService.CreateClassTask(context.Background(), class.OwnerId, task)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)

	res, err = cst.classService.CreateClassTask(context.Background(), xid.New().String(), task)
	assert.NotNil(t, err)

	testCases := []struct{ task domain.ClassTask }{
		{domain.ClassTask{Name: "abelia narindi agsya - abel"}},
		{domain.ClassTask{Description: strings.Repeat("abelia narindi agsya", 520)}},
	}

	for _, tc := range testCases {
		_, err := cst.classService.CreateClassTask(context.Background(), class.OwnerId, &tc.task)
		assert.NotNilf(t, err, "unexpected at %#+v", tc.task)
	}
}

func (cst classServiceTest) listMemberTask(t *testing.T) {
	t.Parallel()

	class := &domain.Class{
		OwnerId: uuid.NewString(),
	}
	err := cst.classService.ClassRepository.CreateClass(context.Background(), class)
	assert.Nil(t, err)

	foo := &domain.User{UserId: uuid.NewString(), Username: "foo", Email: "foo"}
	err = cst.classService.UserRepository.CreateUser(context.Background(), foo)
	assert.Nil(t, err)
	res, err := cst.classService.AddMemberByUsername(context.Background(), class.OwnerId, "foo", &domain.ClassMember{ClassId: class.ClassId})
	if assert.Nil(t, err) {
		assert.Equal(t, 204, res.Code)
	}

	bar := &domain.User{UserId: uuid.NewString(), Username: "bar", Email: "bar"}
	err = cst.classService.UserRepository.CreateUser(context.Background(), bar)
	assert.Nil(t, err)
	res, err = cst.classService.AddMemberByUsername(context.Background(), class.OwnerId, "bar", &domain.ClassMember{ClassId: class.ClassId})
	if assert.Nil(t, err) {
		assert.Equal(t, 204, res.Code)
	}

	for i := 0; i < 10; i++ {
		res, err := cst.classService.AddMember(context.Background(), class.OwnerId, &domain.ClassMember{
			UserId:  uuid.NewString(),
			ClassId: class.ClassId,
		})
		assert.Nil(t, err)
		assert.Equal(t, 204, res.Code)
	}

	resMember, err := cst.classService.ListMember(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, resMember.Code)
	assert.Equal(t, 12, len(resMember.Data))

	_, err = cst.classService.DeleteMember(context.Background(), class.OwnerId, class.ClassId, foo.UserId)

	resMember, err = cst.classService.ListMember(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, resMember.Code)
	assert.Equal(t, 11, len(resMember.Data))
	assert.Nil(t, err)
}
