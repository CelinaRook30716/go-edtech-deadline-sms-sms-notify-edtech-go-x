package main

import (
	"testing"
	"time"
)

func TestShouldAlert(t *testing.T) {
	now := time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		due  time.Time
		want bool
	}{
		{"tomorrow", now.Add(23 * time.Hour), true},
		{"after window", now.Add(25 * time.Hour), false},
		{"past", now.Add(-time.Hour), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ShouldAlert(now, tc.due); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
