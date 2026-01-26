package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/briandowns/spinner"

	"github.com/LeoRiether/gh-box/config"
	"github.com/LeoRiether/gh-box/gh"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Day = 24 * time.Hour

var spin = NewSpinner()

var (
	boxName = flag.String("box", "", "?")
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	flag.Parse()

	cfg := try(config.Get())("getting config")
	box := cfg.Boxes[*boxName]

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
