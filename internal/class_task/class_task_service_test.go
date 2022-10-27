package classtask_test

import (
	"context"
	"strings"
	"testing"

	"nory/domain"
	"nory/internal/class"
	. "nory/internal/class_task"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestClassTaskService(t *testing.T) {
	t.Parallel()
	cts := &ClassTaskService{
		ClassRepository:     class.NewClassRepositoryMem(),
		ClassTaskRepository: NewClassTaskRepositoryMem(),
	}

	t.Run("create get and delete", func(t *testing.T) {
		t.Parallel()

		class := &domain.Class{
			OwnerId: xid.New().String(),
		}

		err := cts.ClassRepository.CreateClass(context.Background(), class)
		assert.Nil(t, err)

		task := &domain.ClassTask{
			ClassId:  class.ClassId,
			AuthorId: class.OwnerId,
		}

		create, err := cts.CreateTask(context.Background(), task)
		assert.Nil(t, err)
		assert.Equal(t, 200, create.Code)

		get, err := cts.GetTask(context.Background(), create.Data.TaskId)
		assert.Nil(t, err)
		assert.Equal(t, 200, get.Code)
		assert.Equal(t, get.Data, create.Data)

		del, err := cts.DeleteTask(context.Background(), create.Data.TaskId, xid.New().String())
		assert.NotNil(t, err)

		del, err = cts.DeleteTask(context.Background(), create.Data.TaskId, task.AuthorId)
		assert.Nil(t, err)
		assert.Equal(t, 204, del.Code)

		del, err = cts.DeleteTask(context.Background(), create.Data.TaskId, task.AuthorId)
		assert.NotNil(t, err)

		get, err = cts.GetTask(context.Background(), create.Data.TaskId)
		assert.NotNil(t, err)

		for _, tc := range []struct {
			Task domain.ClassTask
		}{
			{domain.ClassTask{Name: strings.Repeat("foo", 20)}},
			{domain.ClassTask{Description: strings.Repeat("foo", 1024)}},
		} {
			_, err := cts.CreateTask(context.Background(), &tc.Task)
			assert.NotNil(t, err)
		}
	})
}
