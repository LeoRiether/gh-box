package gh

import (
	"slices"
	"testing"
	"time"

	"github.com/LeoRiether/gh-box/test/assert"
)

func TestMakeSearchPRArgs(t *testing.T) {
	created := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	opts := SearchOptions{
		Authors:      []string{"joaovaladares", "LeoRiether"},
		Organization: "acme",
		CreatedAfter: &created,
		UpdatedAfter: &updated,
		State:        Merged,
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

	assert.Equal(t, got, want)
}

func TestMakeSearchPRArgsMinimal(t *testing.T) {
	opts := SearchOptions{
		State: AnyState,
	}

	got := makeSearchPRArgs(opts)
	want := []string{
		"search", "prs",
		"--json", "author,createdAt,updatedAt,title,state,isDraft,url",
		"--",
	}

	assert.Equal(t, got, want)
}

func TestMakeSearchPRArgsStates(t *testing.T) {
	tests := []struct {
		name  string
		state PRState
		want  string
	}{
		{name: "open", state: Open, want: "state:open"},
		{name: "closed", state: Closed, want: "state:closed"},
		{name: "merged", state: Merged, want: "is:merged"},
		{name: "any", state: AnyState, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeSearchPRArgs(SearchOptions{State: tt.state})
			if tt.want == "" {
				for _, arg := range got {
					assert.False(t, arg == "state:open" || arg == "state:closed" || arg == "is:merged")
				}
				return
			}
			assert.True(t, slices.Contains(got, tt.want))
		})
	}
}
