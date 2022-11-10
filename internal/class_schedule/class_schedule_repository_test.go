package classschedule_test

import (
	"context"
	"os"
	"testing"
	"time"

	"nory/domain"
	"nory/internal/class"
	. "nory/internal/class_schedule"
	"nory/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

var now = time.Now().UTC().Round(time.Hour)

func TestClassScheduleRepository(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}

	repos := []Repository{
		{
			Name:                    "memory",
			ClassScheduleRepository: NewClassScheduleRepositoryMem(),
			ClassRepository:         class.NewClassRepositoryMem(),
			UserRepository:          user.NewUserRepositoryMem(),
		},
		{
			Skip:                    os.Getenv("DATABASE_URL") == "",
			Name:                    "postgres",
			ClassScheduleRepository: NewClassScheduleRepositoryPg(pool),
			ClassRepository:         class.NewClassRepositoryPostgres(pool),
			UserRepository:          user.NewUserRepositoryPostgres(pool),
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
			t.Run("CreateSchedule", repo.testCreateSchedule)
			t.Run("GetSchedules", repo.testGetSchedules)
			t.Run("GetSchedule", repo.testGetSchedule)
			t.Run("ClearSchedules", repo.testClearSchedules)
			t.Run("DeleteSchedule", repo.testDeleteSchedule)
		})
	}
}

type Repository struct {
	Name                    string
	ClassScheduleRepository domain.ClassScheduleRepository
	ClassRepository         domain.ClassRepository
	UserRepository          domain.UserRepository
	Skip                    bool
	Schedules               []domain.ClassSchedule

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

func (r *Repository) testCreateSchedule(t *testing.T) {
	testCases := []struct {
		Name     string
		Schedule domain.ClassSchedule
		Err      error
	}{
		{"success", domain.ClassSchedule{ClassId: r.getClass("classFoo"), AuthorId: r.getUser("classFoo"), Day: int8(2), StartAt: now}, nil},
		{"success", domain.ClassSchedule{ClassId: r.getClass("classFoo"), AuthorId: r.getUser("classFoo"), Day: int8(1), StartAt: now}, nil},
		{"success", domain.ClassSchedule{ClassId: r.getClass("classFoo"), AuthorId: r.getUser("classFoo"), Day: int8(1), StartAt: now}, nil},
		{"success", domain.ClassSchedule{ClassId: r.getClass("classBar"), AuthorId: r.getUser("classBar"), Day: int8(1), StartAt: now}, nil},
		{"success", domain.ClassSchedule{ClassId: r.getClass("classBar"), AuthorId: r.getUser("classBar"), Day: int8(1), StartAt: now}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			s := tc.Schedule
			err := r.ClassScheduleRepository.CreateSchedule(context.Background(), &s)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.NotEqual(t, tc.Schedule.ScheduleId, s.ScheduleId, "CreateSchedule must assign generated id to (ClassSchedule).ScheduleId")
				r.Schedules = append(r.Schedules, s)
			}
		})
	}
}

func (r *Repository) testGetSchedule(t *testing.T) {
	testCases := []struct {
		Name       string
		ScheduleId string
		Err        error
	}{
		{"success", r.Schedules[0].ScheduleId, nil},
		{"not found", xid.New().String(), domain.ErrClassScheduleNotExists},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			schedule, err := r.ClassScheduleRepository.GetSchedule(context.Background(), tc.ScheduleId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			if err == nil {
				assert.Equal(t, tc.ScheduleId, schedule.ScheduleId, "unknown ScheduleId")
			}
		})
	}
}

func (r *Repository) testGetSchedules(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Len     int
		Err     error
	}{
		{"2 task same day 1 different day", r.getClass("classFoo"), 3, nil},
		{"2 task same day", r.getClass("classBar"), 2, nil},
		{"has no task", r.getClass("classBaz"), 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			schedules, err := r.ClassScheduleRepository.GetSchedules(context.Background(), tc.ClassId)
			assert.Equal(t, tc.Err, err, "missmatch error")
			assert.Equal(t, tc.Len, len(schedules), "missmatch schedules length")
			for _, schedule := range schedules {
				assert.Equal(t, schedule.ClassId, tc.ClassId, "unknown ClassId received")
			}
		})
	}
}

func (r *Repository) testClearSchedules(t *testing.T) {
	testCases := []struct {
		Name    string
		ClassId string
		Day     int8
		Len     int
		Err     error
	}{
		{"success", r.getClass("classFoo"), 2, 2, nil},
		{"success", r.getClass("classBar"), 1, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := r.ClassScheduleRepository.ClearSchedules(context.Background(), tc.ClassId, tc.Day)
			assert.Equal(t, tc.Err, err, "missmatch error")

			schedules, err := r.ClassScheduleRepository.GetSchedules(context.Background(), tc.ClassId)
			assert.Nil(t, err)
			assert.Equal(t, tc.Len, len(schedules))
		})
	}
}

func (r *Repository) testDeleteSchedule(t *testing.T) {
	for _, s := range r.Schedules {
		err := r.ClassScheduleRepository.DeleteSchedule(context.Background(), s.ScheduleId)
		assert.Nil(t, err, "unknown error")
		_, err = r.ClassScheduleRepository.GetSchedule(context.Background(), s.ScheduleId)
		assert.Equal(t, domain.ErrClassScheduleNotExists, err, "schedule not deleted properly")
	}
}
