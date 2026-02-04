package main

import "testing"

func TestParseSinceDays(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "days suffix", input: "14d", want: 14},
		{name: "weeks suffix", input: "2w", want: 14},
		{name: "plain number", input: "7", want: 7},
		{name: "zero string", input: "0", want: 0},
		{name: "trimmed", input: " 3d ", want: 3},
		{name: "empty", input: "", want: 0},
		{name: "invalid suffix", input: "3x", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "nonnumeric", input: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSinceDays(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDaysAgo(t *testing.T) {
	if got := daysAgo(0); got != nil {
		t.Fatalf("expected nil for 0 days, got %v", got)
	}
	if got := daysAgo(-1); got != nil {
		t.Fatalf("expected nil for negative days, got %v", got)
	}

	got := daysAgo(1)
	if got == nil {
		t.Fatalf("expected time for 1 day, got nil")
	}
}
