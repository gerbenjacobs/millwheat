package game

import (
	"testing"
	"time"
)

func TestJob_Progress(t *testing.T) {
	type fields struct {
		Hours     int
		Completed time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "3 hour job checked at 2 hours",
			fields: fields{
				Hours:     3,
				Completed: time.Now().Add(1 * time.Hour),
			},
			want: 66,
		},
		{
			name: "3 hour job checked at 3 hours",
			fields: fields{
				Hours:     3,
				Completed: time.Now(),
			},
			want: 100,
		},
		{
			name: "3 hour job checked at 4 hours",
			fields: fields{
				Hours:     3,
				Completed: time.Now().Add(-1 * time.Hour),
			},
			want: 100,
		},
		{
			name: "1 hour job checked at 59 minutes",
			fields: fields{
				Hours:     1,
				Completed: time.Now().Add(1 * time.Minute),
			},
			want: 98, // because of math.Floor
		},
		{
			name: "1 hour job checked straight away",
			fields: fields{
				Hours:     1,
				Completed: time.Now().Add(60 * time.Minute),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				InputJob: InputJob{
					Duration: time.Duration(tt.fields.Hours) * time.Hour,
				},
				Completed: tt.fields.Completed,
				Status:    JobStatusActive,
			}
			if got := j.Progress(); got != tt.want {
				t.Errorf("Progress() = %v, want %v", got, tt.want)
			}
		})
	}
}
