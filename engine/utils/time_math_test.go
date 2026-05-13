package utils

import (
	"testing"
	"time"
)

func TestAddMonths(t *testing.T) {
	tests := []struct {
		name   string
		start  time.Time
		months int
		want   time.Time
	}{
		{
			name:   "Case A: 2024-01-31 + 1 month (Leap year)",
			start:  time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			months: 1,
			want:   time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "Case B: 2023-01-31 + 1 month (Non-leap)",
			start:  time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			months: 1,
			want:   time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "Case C: 2024-08-31 + 1 month",
			start:  time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC),
			months: 1,
			want:   time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "Case D: 2024-03-31 - 1 month",
			start:  time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			months: -1,
			want:   time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddMonths(tt.start, tt.months)
			if !got.Equal(tt.want) {
				t.Errorf("AddMonths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddYears(t *testing.T) {
	tests := []struct {
		name  string
		start time.Time
		years int
		want  time.Time
	}{
		{
			name:  "Leap Year Anniversary: 2024-02-29 + 1 year",
			start: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			years: 1,
			want:  time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "Leap Year Anniversary: 2024-02-29 + 4 years",
			start: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			years: 4,
			want:  time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddYears(tt.start, tt.years)
			if !got.Equal(tt.want) {
				t.Errorf("AddYears() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAgeAtDate(t *testing.T) {
	tests := []struct {
		name       string
		dob        time.Time
		target     time.Time
		wantMonths int
		wantWeeks  int
		wantYears  int
	}{
		{
			name:       "Case A (Boundary): 2024-01-15 to 2024-03-14",
			dob:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			target:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
			wantMonths: 1,
			wantWeeks:  8, // (16 + 29 + 14) = 59 days -> 8 weeks
			wantYears:  0,
		},
		{
			name:       "Case B (Exact): 2024-01-15 to 2024-03-15",
			dob:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			target:     time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			wantMonths: 2,
			wantWeeks:  8, // (16 + 29 + 15) = 60 days -> 8 weeks
			wantYears:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMonths, gotWeeks, gotYears := GetAgeAtDate(tt.dob, tt.target)
			if gotMonths != tt.wantMonths || gotWeeks != tt.wantWeeks || gotYears != tt.wantYears {
				t.Errorf("%s: GetAgeAtDate() = (%v, %v, %v), want (%v, %v, %v)", 
					tt.name, gotMonths, gotWeeks, gotYears, tt.wantMonths, tt.wantWeeks, tt.wantYears)
			}
		})
	}
}

func TestNormalizeVaccineName(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "Case from Phase 02",
			raw:  "  Vắc-xin   6 trong 1  ",
			want: "vac-xin 6 trong 1",
		},
		{
			name: "Lao with accents and parens",
			raw:  "Vắc xin phòng Lao (BCG)",
			want: "vac xin phong lao",
		},
		{
			name: "Đ character and multiple spaces",
			raw:  "Bạch hầu -   Đuốn ván",
			want: "bach hau - duon van",
		},
		{
			name: "Suffix removal (ml)",
			raw:  "Infanrix Hexa 0.5ml",
			want: "infanrix hexa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeVaccineName(tt.raw)
			if got != tt.want {
				t.Errorf("%s: NormalizeVaccineName() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
