package time

import (
	"testing"
	"time"
)

func TestParseTimeString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   time.Time
		wantOk bool
		isNow  bool // Set to true if the output should be time.Now()
	}{
		{
			name:   "Unix timestamp",
			input:  "1718700900", // Example: 2024-06-18 11:35:00 UTC
			want:   time.Unix(1718700900, 0),
			wantOk: true,
		},
		{
			name:   "DD/MM/YYYY date",
			input:  "18/06/2025",
			want:   time.Date(2025, time.June, 18, 0, 0, 0, 0, time.UTC),
			wantOk: true,
		},
		{
			name:   "DD/MM/YYYY HH:MM:SS",
			input:  "18/06/2025 11:35:00",
			want:   time.Date(2025, time.June, 18, 11, 35, 0, 0, time.UTC),
			wantOk: true,
		},
		{
			name:   "RFC3339",
			input:  "2025-06-18T11:35:00Z",
			want:   time.Date(2025, time.June, 18, 11, 35, 0, 0, time.UTC),
			wantOk: true,
		},
		{
			name:   "Invalid American date (MM/DD/YYYY)",
			input:  "06/18/2025",
			want:   time.Time{},
			wantOk: false,
		},
		{
			name:   "Invalid format",
			input:  "not a date",
			want:   time.Time{},
			wantOk: false,
		},
		{
			name:   "Empty input returns now",
			input:  "",
			want:   time.Time{}, // Ignored, see code below
			wantOk: true,
			isNow:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, _ := ParseTimeString(tt.input)
			if ok != tt.wantOk {
				t.Errorf("ParseTimeString(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
				return
			}
			if !tt.wantOk {
				return // Skip time comparison if parsing should fail
			}
			if tt.isNow {
				// Check that the returned time is within 1 second of now
				now := time.Now()
				diff := now.Sub(got)
				if diff < -time.Second || diff > time.Second {
					t.Errorf("ParseTimeString(%q) = %v, not within 1 second of now (%v)", tt.input, got, now)
				}
			} else if !got.Equal(tt.want) {
				t.Errorf("ParseTimeString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
