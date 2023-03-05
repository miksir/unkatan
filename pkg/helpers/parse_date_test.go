package helpers

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	type args struct {
		str     string
		timeNow time.Time
	}
	timeNow := time.Date(2020, time.Month(1), 2, 19, 37, 0, 0, time.Now().Location())
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "DD-MM-YYYY HH:MM",
			args:    args{"Closed till 02-02-2023 18:53", timeNow},
			want:    time.Date(2023, time.Month(2), 2, 18, 53, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "HH:MM",
			args:    args{"Closed till 15:22", timeNow},
			want:    time.Date(2020, time.Month(1), 2, 15, 22, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "DD.MM HH:MM 1",
			args:    args{"Closed till 15/05 16:22", timeNow},
			want:    time.Date(2020, time.Month(5), 15, 16, 22, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "DD.MM HH:MM 2",
			args:    args{"Closed till 15.05 16:22", timeNow},
			want:    time.Date(2020, time.Month(5), 15, 16, 22, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "YYYY-DD-MM HH:MM",
			args:    args{"Closed till 2024-05-31 18:55", timeNow},
			want:    time.Date(2024, time.Month(5), 31, 18, 55, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "tomorrow HH:MM",
			args:    args{"till tomorrow 18:55", timeNow},
			want:    time.Date(2020, time.Month(1), 3, 18, 55, 0, 0, timeNow.Location()),
			wantErr: false,
		},
		{
			name:    "no till",
			args:    args{"Closed forever", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "no pattern",
			args:    args{"Closed till 01-01-01 01-01", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "Ignore wrong month",
			args:    args{"Closed till 02-13-2023 18:53", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "Ignore wrong day",
			args:    args{"Closed till 32-12-2023 18:53", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "Ignore wrong year",
			args:    args{"Closed till 12-12-2019 18:53", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "Ignore wrong hour",
			args:    args{"Closed till 12-12-2044 24:00", timeNow},
			want:    timeNow,
			wantErr: true,
		},
		{
			name:    "Ignore wrong min",
			args:    args{"Closed till 12-12-2044 22:61", timeNow},
			want:    timeNow,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, endPos, err := ParseDate("till", tt.args.str, tt.args.timeNow)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() got = %v, want %v", got, tt.want)
			}
			if !tt.wantErr && endPos != len(tt.args.str) {
				t.Errorf("ParseDate() end position got %d, want %d", endPos, len(tt.args.str))
			}
		})
	}
}
