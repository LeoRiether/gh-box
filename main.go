package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	"github.com/briandowns/spinner"

	"github.com/LeoRiether/gh-box/config"
	"github.com/LeoRiether/gh-box/gh"
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
	Box string `arg:"" optional:""`
}

func (b *BoxCmd) Run() error {
	cfg := try(config.Load())("getting config")
	box := try(cfg.Box(b.Box))(fmt.Sprintf("box=%q", b.Box))

	spin.Start()

	spin.Message("searching PRs")
	prs := try(gh.SearchPRs(box.People, box.Organization, time.Now().Add(-14*Day)))("searching PRs")

	spin.Message("fetching PR details")
	prdetails := try(gh.ViewPRsDetails(prs))("fetching PR details")

	slices.SortFunc(prdetails, func(a, b gh.PRDetails) int {
		return -a.UpdatedAt.Compare(b.UpdatedAt)
	})

	spin.Stop()
	fmt.Println(prdetails.Style())
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
