package class_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"nory/domain"
	. "nory/internal/class"
	classmember "nory/internal/class_member"
	classschedule "nory/internal/class_schedule"
	classtask "nory/internal/class_task"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestClassService(t *testing.T) {
	t.Parallel()
	classService := ClassService{
		UserRepository:          user.NewUserRepositoryMem(),
		ClassRepository:         NewClassRepositoryMem(),
		ClassTaskRepository:     classtask.NewClassTaskRepositoryMem(),
		ClassMemberRepository:   classmember.NewClassMemberRepositoryMem(),
		ClassScheduleRepository: classschedule.NewClassScheduleRepositoryMem(),
	}

	cst := classServiceTest{classService}

	t.Run("get class info", cst.testClassInfo)
	t.Run("get class tasks", cst.testClassTasks)
	t.Run("create class tasks", cst.testCreateClassTask)
	t.Run("create class Schedule", cst.testClassSchedule)
	t.Run("create, access and delete class", cst.testClassCreate)
	t.Run("list member", cst.testListMember)
}

type classServiceTest struct {
	classService ClassService
}

func (cst classServiceTest) testClassInfo(t *testing.T) {
	t.Parallel()

	_, err := cst.classService.GetClassInfo(context.Background(), "foobarbazqux")
	assert.ErrorContains(t, err, "can not find class with id \"foobarbazqux\"")
	u := &domain.User{
		UserId: uuid.NewString(),
		Name: xid.New().String(),
		Username: xid.New().String(),
		Email: xid.New().String(),
	}

	err = cst.classService.UserRepository.CreateUser(context.Background(), u)
	assert.Nil(t, err)

	classA := &domain.Class{OwnerId: u.UserId, Name: xid.New().String()}
	_, err = cst.classService.CreateClass(context.Background(), classA)
	assert.Nil(t, err)

	res, err := cst.classService.GetClassInfo(context.Background(), classA.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, classA, res.Data)
	res, err = cst.classService.GetClassInfoByName(context.Background(), u.Username, classA.Name)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, classA, res.Data)
}

func (cst classServiceTest) testClassTasks(t *testing.T) {
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

	res, err := cst.classService.GetClassTasks(context.Background(), class.ClassId, time.Time{}, time.Time{})
	assert.Nil(t, err)

	for _, task := range res.Data {
		_, err := cst.classService.DeleteClassTask(context.Background(), class.OwnerId, task.TaskId)
		assert.Nil(t, err)
	}

	res, err = cst.classService.GetClassTasks(context.Background(), class.ClassId, time.Time{}, time.Time{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(res.Data))
}

func (cst classServiceTest) testClassSchedule(t *testing.T) {
	t.Parallel()

	class := &domain.Class{
		OwnerId: uuid.NewString(),
		Name:    "foo",
	}

	r, err := cst.classService.CreateClass(context.Background(), class)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.Code)

	_, err = cst.classService.GetSchedule(context.Background(), xid.New().String())
	assert.NotNil(t, err)

	for i := 0; i < 7; i++ {
		schedule := &domain.ClassSchedule{
			AuthorId: class.OwnerId,
			ClassId:  class.ClassId,
			Name:     "MATH!!!",
			Day:      int8(i),
			StartAt:  time.Now().UTC().Round(time.Hour),
			Duration: int16(20),
		}
		_, err := cst.classService.CreateSchedule(context.Background(), schedule)
		assert.Nil(t, err)

		res, err := cst.classService.GetSchedule(context.Background(), schedule.ScheduleId)
		res.Data.CreatedAt = time.Time{}
		assert.Nil(t, err)
		assert.Equal(t, schedule, res.Data)

		t.Cleanup(func() {
			_, err := cst.classService.DeleteSchedule(context.Background(), uuid.NewString(), schedule.ScheduleId)
			assert.NotNil(t, err)
			_, err = cst.classService.DeleteSchedule(context.Background(), class.OwnerId, schedule.ScheduleId)
			assert.Nil(t, err)
		})
	}

	schedules, err := cst.classService.GetClassSchedules(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(schedules.Data))

	{

		class := &domain.Class{
			OwnerId: uuid.NewString(),
			Name:    "foo",
		}

		r, err = cst.classService.CreateClass(context.Background(), class)
		assert.Nil(t, err)
		assert.Equal(t, 200, r.Code)

		for i := 0; i < 7; i++ {
			schedule := &domain.ClassSchedule{
				AuthorId: class.OwnerId,
				ClassId:  class.ClassId,
				Name:     "MATH!!!",
				Day:      int8(i),
				StartAt:  time.Now().UTC().Round(time.Hour),
				Duration: int16(20),
			}
			_, err := cst.classService.CreateSchedule(context.Background(), schedule)
			assert.Nil(t, err)
		}

		for i := 0; i < 7; i++ {
			schedules, err := cst.classService.GetClassSchedules(context.Background(), class.ClassId)
			assert.Nil(t, err)
			assert.Equal(t, 7-i, len(schedules.Data))

			_, err = cst.classService.ClearSchedules(context.Background(), uuid.NewString(), class.ClassId, int8(i))
			assert.NotNil(t, err)
			_, err = cst.classService.ClearSchedules(context.Background(), class.OwnerId, class.ClassId, int8(i))
			assert.Nil(t, err)

			schedules, err = cst.classService.GetClassSchedules(context.Background(), class.ClassId)
			assert.Nil(t, err)
			assert.Equal(t, 6-i, len(schedules.Data))
		}
	}
}

func (cst classServiceTest) testClassCreate(t *testing.T) {
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

	_, err = cst.classService.UpdateClass(context.Background(), admin, &domain.Class{
		Name: "foo",
		ClassId: class.ClassId,
	})
	assert.NotNil(t, err)

	_, err = cst.classService.UpdateClass(context.Background(), class.OwnerId, &domain.Class{
		Name: "foo",
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

	err = cst.classService.AccessClass(context.Background(), uuid.NewString(), class.ClassId, "admin")
	assert.NotNil(t, err)

	err = cst.classService.AccessClass(context.Background(), class.OwnerId, class.ClassId, "admin")
	assert.Nil(t, err)

	err = cst.classService.AccessClass(context.Background(), member, class.ClassId, "admin")
	assert.NotNil(t, err)

	err = cst.classService.AccessClass(context.Background(), admin, class.ClassId, "admin")
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

		_, err = cst.classService.DeleteClass(context.Background(), uuid.NewString(), class.ClassId)
		assert.NotNil(t, err)
	}
}

func (cst *classServiceTest) testCreateClassTask(t *testing.T) {
	t.Parallel()
	class := &domain.Class{
		Name: "foo",
		OwnerId: uuid.NewString(),
	}
	_, err := cst.classService.CreateClass(context.Background(), class)
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

func (cst classServiceTest) testListMember(t *testing.T) {
	t.Parallel()

	class := &domain.Class{
		OwnerId: uuid.NewString(),
		Name: "foo",
	}
	_, err := cst.classService.CreateClass(context.Background(), class)
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
	assert.Equal(t, 13, len(resMember.Data))

	_, err = cst.classService.DeleteMember(context.Background(), class.OwnerId, class.ClassId, foo.UserId)

	resMember, err = cst.classService.ListMember(context.Background(), class.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 200, resMember.Code)
	assert.Equal(t, 12, len(resMember.Data))
	assert.Nil(t, err)

	_, err = cst.classService.UpdateMember(context.Background(), class.OwnerId, &domain.ClassMember{
		ClassId: class.ClassId,
		UserId:  bar.UserId,
		Level:   "admin",
	})
	assert.Nil(t, err)

	resMember, err = cst.classService.ListMember(context.Background(), class.ClassId)
	assert.Nil(t, err)
	for _, i := range resMember.Data {
		if i.UserId == bar.UserId {
			assert.Equal(t, "admin", i.Level)
		}
	}
}
