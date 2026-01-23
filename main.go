package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/LeoRiether/gh-box/gh"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Day = 24 * time.Hour

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	prs := try(gh.SearchPRs(
		[]string{"LeoRiether", "qrno", "joaovaladares", "figueredo",
			"fabricio-suarte", "daviromao", "gabrielpessoa1"},
		time.Now().Add(-14*Day),
	))

	slices.SortFunc(prs, func(a, b gh.PullRequest) int { return a.UpdatedAt.Compare(b.UpdatedAt) })

	for _, pr := range prs {
		fmt.Printf("[%s] %s -- by %s\n",
			pr.State, pr.Title, pr.Author.Login)
	}
}

func try[T any](value T, err error) T {
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return value
}
