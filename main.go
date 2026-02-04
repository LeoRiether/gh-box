package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	"github.com/briandowns/spinner"

	"github.com/LeoRiether/gh-box/config"
	"github.com/LeoRiether/gh-box/gh"
	"github.com/LeoRiether/gh-box/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Day = 24 * time.Hour

var spin = NewSpinner()

type CLI struct {
	Box        BoxCmd               `cmd:"" default:"withargs" help:"Choose a PR box"`
	ConfigPath config.ConfigPathCmd `cmd:"" help:"Shows the configuration path"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cli := CLI{}
	_ = kong.Parse(&cli).Run()
}

type BoxCmd struct {
	Box          string   `arg:"" optional:""`
	Authors      []string `name:"authors" help:"Comma-separated list of authors; overrides box people" sep:","`
	State        string   `help:"Filter by PR state (all, open, closed, merged)" default:"all" enum:"all,open,closed,merged"`
	CreatedSince string   `name:"created-since" help:"Only PRs created in the last N days (e.g. 14d, 2w)" default:"14d"`
	UpdatedSince string   `name:"updated-since" help:"Only PRs updated in the last N days (e.g. 7d, 2w)" default:""`
}

func (b *BoxCmd) Run() error {
	cfg := try(config.Load())("getting config")
	box := try(cfg.Box(b.Box))(fmt.Sprintf("box=%q", b.Box))

	authors := box.People
	if len(b.Authors) > 0 {
		authors = b.Authors
	}

	createdSince := try(util.ParseDuration(b.CreatedSince))("--created-since")
	updatedSince := try(util.ParseDuration(b.UpdatedSince))("--updated-since")

	opts := gh.SearchOptions{
		Authors:      authors,
		Organization: box.Organization,
		CreatedAfter: createdSince.Ago(),
		UpdatedAfter: updatedSince.Ago(),
		State:        gh.PRState(b.State),
	}

	spin.Start()

	spin.Message("searching PRs")
	prs := try(gh.SearchPRs(opts))("searching PRs")

	spin.Message("fetching PR details")
	prdetails := try(gh.ViewPRsDetails(prs))("fetching PR details")

	slices.SortFunc(prdetails, func(a, b gh.PRDetails) int {
		return -a.UpdatedAt.Compare(b.UpdatedAt)
	})

	spin.Stop()
	pager(prdetails.Style())
	return nil
}

type Spinner struct {
	*spinner.Spinner
}

func NewSpinner() *Spinner {
	return &Spinner{
		Spinner: spinner.New(
			spinner.CharSets[14],
			100*time.Millisecond,
			spinner.WithWriter(os.Stderr)),
	}
}

func (s *Spinner) Message(message string) {
	s.Suffix = " " + message
}

func try[T any](value T, err error) func(message string) T {
	return func(message string) T {
		if err != nil {
			spin.Stop()
			fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
			os.Exit(1)
		}

		return value
	}
}

// Calls `less -r` with the content passed in
func pager(content string) error {
	cmd := exec.Command("less", "--raw-control-chars", "--quit-if-one-screen", "--no-init")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	fmt.Fprint(stdin, content)
	stdin.Close()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
