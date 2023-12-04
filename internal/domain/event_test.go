package domain_test

import (
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	a := assert.New(t)
	type args struct {
		name        string
		description *string
		place       string
		ageRating   int
		startAt     time.Time
		endAt       time.Time
	}
	testDesc := "test desc"
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without desc",
			args: args{
				name:        "Test",
				description: nil,
				place:       "test place",
				ageRating:   16,
				startAt:     time.Now().UTC().Add(1 * time.Hour),
				endAt:       time.Now().UTC().Add(2 * time.Hour),
			},
		},
		{
			name: "with desc",
			args: args{
				name:        "Test",
				description: &testDesc,
				place:       "test place",
				ageRating:   16,
				startAt:     time.Now().UTC().Add(1 * time.Hour),
				endAt:       time.Now().UTC().Add(2 * time.Hour),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := domain.NewEvent(
				tt.args.name,
				tt.args.description,
				tt.args.place,
				tt.args.ageRating,
				tt.args.startAt,
				tt.args.endAt,
			)
			a.NotEqual(uuid.Nil, e.ID)
			a.Equal(tt.args.name, e.Name)
			if tt.args.description != nil {
				a.Equal(*tt.args.description, *e.Description)
			} else {
				a.Nil(e.Description)
			}
			a.Equal(tt.args.place, e.Place)
			a.Equal(tt.args.ageRating, e.AgeRating)
			a.Equal(tt.args.startAt, e.StartAt)
			a.Equal(tt.args.endAt, e.EndAt)
			a.False(e.CreatedAt.IsZero())
			a.False(e.UpdatedAt.IsZero())
			a.WithinDuration(e.CreatedAt, e.UpdatedAt, time.Millisecond)
		})
	}
}
