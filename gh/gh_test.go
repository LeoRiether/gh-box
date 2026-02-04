package gh

import (
	"strings"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMakeSearchPRArgs(t *testing.T) {
	run := func(opts SearchOptions) string {
		args := makeSearchPRArgs(opts)
		snap := strings.Join(args, " ")
		snap = strings.Replace(snap, "author,createdAt,updatedAt,title,state,isDraft,url", "<json>", 1)
		return snap
	}

	t.Run("Standard", func(t *testing.T) {
		created := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
		updated := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		snap := run(SearchOptions{
			Authors:      []string{"joaovaladares", "LeoRiether"},
			Organization: "acme",
			CreatedAfter: &created,
			UpdatedAfter: &updated,
			State:        Merged,
		})

		snaps.MatchInlineSnapshot(t, snap, snaps.Inline("search prs --json <json> -- --created >=2026-01-02 --updated >=2026-01-10 (author:joaovaladares OR author:LeoRiether) org:acme is:merged"))
	})

	t.Run("Single author", func(t *testing.T) {
		snap := run(SearchOptions{
			Authors: []string{"@me"},
		})

		snaps.MatchInlineSnapshot(t, snap, snaps.Inline("search prs --json <json> -- (author:@me)"))
	})

	t.Run("Minimal", func(t *testing.T) {
		snap := run(SearchOptions{
			State: AnyState,
		})

		snaps.MatchInlineSnapshot(t, snap, snaps.Inline("search prs --json <json> --"))
	})

	t.Run("Match specific PR state", func(t *testing.T) {
		snap := run(SearchOptions{
			State: Closed,
		})

		snaps.MatchInlineSnapshot(t, snap, snaps.Inline("search prs --json <json> -- state:closed"))
	})
}
