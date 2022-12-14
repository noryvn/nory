package classtask_test

import (
	"context"
	"os"
	"testing"
	"time"

	"nory/domain"
	"nory/internal/class"
	. "nory/internal/class_task"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

var (
	SpecialDate = time.Date(2005, time.August, 11, 0, 0, 0, 0, time.UTC)
	Now         = time.Now().UTC().Round(24 * time.Hour)
	Tomorrow    = Now.Add(24 * time.Hour)
)

func TestClassTaskRepository(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}

	repos := []Repository{
		{
			Name:                "memory",
			ClassTaskRepository: NewClassTaskRepositoryMem(),
			ClassRepository:     class.NewClassRepositoryMem(),
			UserRepository:      user.NewUserRepositoryMem(),
		},
		{
			Name:                "postgres",
			ClassTaskRepository: NewClassTaskRepositoryPostgres(pool),
			ClassRepository:     class.NewClassRepositoryPostgres(pool),
			UserRepository:      user.NewUserRepositoryPostgres(pool),
			Skip:                os.Getenv("DATABASE_URL") == "",
		},
	}

	for _, repo := range repos {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
			repo.t = t
			if repo.Skip {
				t.Skipf("skipping %s", repo.Name)
			}
			t.Parallel()
			t.Run("CreateTask", repo.testCreateTask)
			t.Run("GetTask", repo.testGetTask)
			t.Run("GetTasks", repo.testGetTasks)
			t.Run("GetTasksWithRange", repo.testGetTasksWithRange)
			t.Run("UpdateTask", repo.testUpdateTasks)
			t.Run("DeleteTask", repo.testDeleteTask)
		})
	}
}

type Repository struct {
	Name                string
	ClassTaskRepository domain.ClassTaskRepository
	ClassRepository     domain.ClassRepository
	UserRepository      domain.UserRepository
	Skip                bool

	tasks   []domain.ClassTask
	classes map[string]*domain.Class
	t       *testing.T
}

func (r *Repository) getClass(name string) string {
	if r.classes == nil {
		r.classes = make(map[string]*domain.Class)
	}
	if class, ok := r.classes[name]; ok {
		return class.ClassId
	}

	u := &domain.User{
		UserId:   uuid.NewString(),
		Email:    xid.New().String(),
		Username: xid.New().String(),
	}
	err := r.UserRepository.CreateUser(context.Background(), u)
	assert.Nil(r.t, err)

	class := &domain.Class{
		Name:    xid.New().String(),
		OwnerId: u.UserId,
	}
	err = r.ClassRepository.CreateClass(context.Background(), class)
	assert.Nil(r.t, err)

	r.classes[name] = class
	return class.ClassId
}

func (r *Repository) getUser(name string) string {
	r.getClass(name)
	return r.classes[name].OwnerId
}

func (r *Repository) testCreateTask(t *testing.T) {
	testCases := []struct {
		Name string
		Task domain.ClassTask
		Err  bool
	}{
		{"success", domain.ClassTask{ClassId: r.getClass("foo"), AuthorId: r.getUser("foo"), DueDate: SpecialDate, TaskId: "abelia narindi agsya"}, false},
		{"success", domain.ClassTask{ClassId: r.getClass("foo"), AuthorId: r.getUser("foo"), DueDate: Now, Name: "abelia narindi agsya"}, false},
		{"success", domain.ClassTask{ClassId: r.getClass("bar"), AuthorId: r.getUser("bar"), DueDate: SpecialDate, Description: "abelia narindi agsya"}, false},
		{"success", domain.ClassTask{ClassId: r.getClass("baz"), AuthorId: r.getUser("baz"), DueDate: Now}, false},
		{"failed empty ClassId", domain.ClassTask{AuthorId: "", DueDate: SpecialDate, Name: "Err"}, true},
		{"failed empty AuthorId", domain.ClassTask{ClassId: r.getClass("foo"), DueDate: SpecialDate, Name: "Err"}, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			task := tc.Task
			err := r.ClassTaskRepository.CreateTask(context.Background(), &task)
			if tc.Err {
				assert.NotNil(t, err, "missmatch error")
			} else {
				assert.Nil(t, err, "unexpected error %v", err)
				assert.NotEqual(t, tc.Task.TaskId, task.TaskId, "CreateTask should update (*ClassTask).TaskId to generated id")
				r.tasks = append(r.tasks, task)
			}
		})
	}
}

func (r *Repository) testGetTask(t *testing.T) {
	for _, taskSc := range r.tasks {
		task, err := r.ClassTaskRepository.GetTask(context.Background(), taskSc.TaskId)
		assert.Nil(t, err, "unexpected error")
		task.CreatedAt = time.Time{}
		assert.Equal(t, taskSc, *task, "unknown task")
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
			id := r.getClass(tc.ClassId)
			tasks, err := r.ClassTaskRepository.GetTasks(context.Background(), id)
			assert.Equal(t, tc.Err, err, "missmatch error")
			assert.Equal(t, tc.Len, len(tasks), "unexpected result length")
			for _, task := range tasks {
				assert.Equal(t, id, task.ClassId, "unknown ClassId")
			}
		})
	}
}

func (r *Repository) testGetTasksWithRange(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		From    time.Time
		To      time.Time
		Len     int
		Err     error
	}{
		{"success", "foo", SpecialDate, Tomorrow, 2, nil},
		{"success", "bar", SpecialDate, Tomorrow, 1, nil},
		{"success", "baz", SpecialDate, Tomorrow, 1, nil},
		{"success", "qux", SpecialDate, Tomorrow, 0, nil},
		{"success", "foo", SpecialDate, SpecialDate.Add(24 * time.Hour), 1, nil},
		{"success", "foo", Now, Tomorrow, 1, nil},
		{"bar has 0 with due date now", "bar", Now, Tomorrow, 0, nil},
		{"baz has 1 with due date Now", "baz", Now, Tomorrow, 1, nil},
		{"ClassID not exists", "qux", Now, Tomorrow, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			id := r.getClass(tc.ClassId)
			tasks, err := r.ClassTaskRepository.GetTasksWithRange(context.Background(), id, tc.From, tc.To)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if !assert.Equal(t, tc.Len, len(tasks), "unexpected result length") {
				for _, task := range tasks {
					t.Logf("%#+v\n", task)
				}
			}
			for _, task := range tasks {
				assert.Equal(t, id, task.ClassId, "unknown ClassId")
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
		{"success", r.tasks[0], nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			before := tc.Task
			before.Name = "Abelia"
			before.Description = xid.New().String()

			err := r.ClassTaskRepository.UpdateTask(context.Background(), &before)
			assert.Equal(t, tc.Err, err)
			if err != nil {
				return
			}

			after, err := r.ClassTaskRepository.GetTask(context.Background(), before.TaskId)
			after.CreatedAt = time.Time{}
			assert.Nil(t, err)
			assert.Equal(t, before, *after)
			assert.NotEqual(t, tc.Task, *after)
		})
	}
}

func (r *Repository) testDeleteTask(t *testing.T) {
	for _, task := range r.tasks {
		err := r.ClassTaskRepository.DeleteTask(context.Background(), task.TaskId)
		assert.Nil(t, err, "unexpected error")

		_, err = r.ClassTaskRepository.GetTask(context.Background(), task.TaskId)
		assert.Equal(t, domain.ErrClassTaskNotExists, err, "failed deleting task")
	}
}
