package gh

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestMakeSearchPRArgs(t *testing.T) {
	created := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	opts := SearchOptions{
		Authors:      []string{"joaovaladares", "LeoRiether"},
		Organization: "acme",
		CreatedAfter: &created,
		UpdatedAfter: &updated,
		State:        StateMerged,
	}

	got := makeSearchPRArgs(opts)
	want := []string{
		"search", "prs",
		"--json", "author,createdAt,updatedAt,title,state,isDraft,url",
		"--",
		"--created", ">=2026-01-02",
		"--updated", ">=2026-01-10",
		"(author:joaovaladares", "OR", "author:LeoRiether)",
		"org:acme",
		"is:merged",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestMakeSearchPRArgsMinimal(t *testing.T) {
	opts := SearchOptions{
		State: StateAny,
	}

	got := makeSearchPRArgs(opts)
	want := []string{
		"search", "prs",
		"--json", "author,createdAt,updatedAt,title,state,isDraft,url",
		"--",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestMakeSearchPRArgsStates(t *testing.T) {
	tests := []struct {
		name  string
		state PRStateFilter
		want  string
	}{
		{name: "open", state: StateOpen, want: "state:open"},
		{name: "closed", state: StateClosed, want: "state:closed"},
		{name: "merged", state: StateMerged, want: "is:merged"},
		{name: "any", state: StateAny, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeSearchPRArgs(SearchOptions{State: tt.state})
			if tt.want == "" {
				for _, arg := range got {
					if arg == "state:open" || arg == "state:closed" || arg == "is:merged" {
						t.Fatalf("unexpected state arg %q in %#v", arg, got)
					}
				}
				return
			}
			if !slices.Contains(got, tt.want) {
				t.Fatalf("missing %q in %#v", tt.want, got)
			}
		})
	}
}
