package class_task_test

import (
	"context"
	"testing"
	"time"

	"nory/domain"
	. "nory/internal/class_task"

	"github.com/stretchr/testify/assert"
)

var (
	AbelBirthday = time.Date(2005, time.August, 11, 0, 0, 0, 0, time.UTC)
	Now          = time.Now()
)

func TestClassTaskRepository(t *testing.T) {
	repos := []Repository{
		{
			Name: "memory",
			R:    NewClassTaskRepositoryMem(),
		},
	}

	for _, repo := range repos {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			if repo.Skip {
				t.Skipf("skipping %s", repo.Name)
			}
			t.Parallel()
			t.Run("CreateTask", repo.testCreateTask)
			t.Run("GetTasks", repo.testGetTask)
			t.Run("GetTasks", repo.testGetTasks)
			t.Run("GetTasksWithDate", repo.testGetTasksWithDate)
			t.Run("UpdateTask", repo.testUpdateTasks)
			t.Run("DeleteTask", repo.testDeleteTask)
		})
	}
}

type Repository struct {
	Name  string
	R     domain.ClassTaskRepository
	Skip  bool
	Tasks []domain.ClassTask
}

func (r *Repository) testCreateTask(t *testing.T) {
	testCases := []struct {
		Name string
		Task domain.ClassTask
		Err  error
	}{
		{"success", domain.ClassTask{ClassId: "foo", DueDate: AbelBirthday, TaskId: "abelia narindi agsya"}, nil},
		{"success", domain.ClassTask{ClassId: "foo", DueDate: Now, Name: "abelia narindi agsya"}, nil},
		{"success", domain.ClassTask{ClassId: "bar", DueDate: AbelBirthday, Description: "abelia narindi agsya"}, nil},
		{"success", domain.ClassTask{ClassId: "baz", DueDate: Now}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			task := tc.Task
			err := r.R.CreateTask(context.Background(), &task)
			assert.Equal(t, tc.Err, err, "missmatch error")
			assert.NotEqual(t, tc.Task.TaskId, task.TaskId, "CreateTask should update (*ClassTask).TaskId to generated id")
			if err == nil {
				r.Tasks = append(r.Tasks, task)
			}
		})
	}
}

func (r *Repository) testGetTask(t *testing.T) {
	for _, taskSc := range r.Tasks {
		task, err := r.R.GetTask(context.Background(), taskSc.TaskId)
		assert.Equal(t, nil, err, "unexpected error")
		assert.Equal(t, taskSc, *task, "unknown TaskId")
	}
}

func (r *Repository) testGetTasks(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Len     int
		Err     error
	}{
		{"success", "foo", 2, nil},
		{"success", "bar", 1, nil},
		{"success", "baz", 1, nil},
		{"success", "qux", 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			tasks, err := r.R.GetTasks(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			assert.Equal(t, tc.Len, len(tasks), "unexpected result length")
			for _, task := range tasks {
				assert.Equal(t, tc.ClassId, task.ClassId, "unknown ClassId")
			}
		})
	}
}

func (r *Repository) testGetTasksWithDate(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		DueDate time.Time
		Len     int
		Err     error
	}{
		{"success", "foo", AbelBirthday, 1, nil},
		{"success", "foo", Now, 1, nil},
		{"bar has 0 with due date now", "bar", Now, 0, nil},
		{"baz has 1 with due date Now", "baz", Now, 1, nil},
		{"ClassID not exists", "qux", Now, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			tasks, err := r.R.GetTasksWithDate(context.Background(), tc.ClassId, tc.DueDate)
			assert.Equal(t, tc.Err, err, "missmatch error")
			assert.Equal(t, tc.Len, len(tasks), "unexpected result length")
			for _, task := range tasks {
				assert.Equal(t, tc.ClassId, task.ClassId, "unknown ClassId")
				assert.Equal(t, tc.DueDate, task.DueDate, "unknown DueDate")
			}
		})
	}
}

func (r *Repository) testUpdateTasks(t *testing.T) {
	testCases := []struct {
		Name string
		Task domain.ClassTask
		Err  error
	}{
		{"success", r.Tasks[0], nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
		})
	}
}

func (r *Repository) testDeleteTask(t *testing.T) {
	for _, task := range r.Tasks {
		err := r.R.DeleteTask(context.Background(), task.TaskId)
		assert.Equal(t, nil, err, "unexpected error")
	}

	testCases := []struct {
		ClassId string
	}{
		{"foo"},
		{"bar"},
		{"baz"},
	}

	for _, tc := range testCases {
		tasks, err := r.R.GetTasks(context.Background(), tc.ClassId)
		assert.Equal(t, nil, err)
		assert.Equal(t, 0, len(tasks), "failed deleting task")
	}
}
