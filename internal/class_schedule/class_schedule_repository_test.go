package classschedule_test

import (
	"context"
	"testing"

	"nory/domain"
	. "nory/internal/class_schedule"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

var (
	classFoo = xid.New().String()
	classBar = xid.New().String()
	classBaz = xid.New().String()
)

func TestClassScheduleRepository(t *testing.T) {
	repos := []Repository{
		{
			Name: "memory",
			R:    NewClassScheduleRepositoryMem(),
		},
	}

	for _, repo := range repos {
		repo := repo
		t.Run(repo.Name, func(t *testing.T) {
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
	Name      string
	R         domain.ClassScheduleRepository
	Skip      bool
	Schedules []domain.ClassSchedule
}

func (r *Repository) testCreateSchedule(t *testing.T) {
	testCases := []struct {
		Name     string
		Schedule domain.ClassSchedule
		Err      error
	}{
		{"success", domain.ClassSchedule{ClassId: classFoo, Day: int8(2)}, nil},
		{"success", domain.ClassSchedule{ClassId: classFoo, Day: int8(1)}, nil},
		{"success", domain.ClassSchedule{ClassId: classFoo, Day: int8(1)}, nil},
		{"success", domain.ClassSchedule{ClassId: classBar, Day: int8(1)}, nil},
		{"success", domain.ClassSchedule{ClassId: classBar, Day: int8(1)}, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			s := tc.Schedule
			err := r.R.CreateSchedule(context.Background(), &s)
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
			schedule, err := r.R.GetSchedule(context.Background(), tc.ScheduleId)
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
		{"2 class same day 1 different day", classFoo, 3, nil},
		{"2 class same day", classBar, 2, nil},
		{"has no class", classBaz, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			schedules, err := r.R.GetSchedules(context.Background(), tc.ClassId)
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
		{"success", classFoo, 2, 2, nil},
		{"success", classBar, 1, 0, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			err := r.R.ClearSchedules(context.Background(), tc.ClassId, tc.Day)
			assert.Equal(t, tc.Err, err, "missmatch error")

			schedules, err := r.R.GetSchedules(context.Background(), tc.ClassId)
			assert.Nil(t, err)
			assert.Equal(t, tc.Len, len(schedules))
		})
	}
}

func (r *Repository) testDeleteSchedule(t *testing.T) {
	for _, s := range r.Schedules {
		err := r.R.DeleteSchedule(context.Background(), s.ScheduleId)
		assert.Nil(t, err, "unknown error")
		_, err = r.R.GetSchedule(context.Background(), s.ScheduleId)
		assert.Equal(t, domain.ErrClassScheduleNotExists, err, "schedule not deleted properly")
	}
}
