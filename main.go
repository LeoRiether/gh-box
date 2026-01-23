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
		"inloco",
		time.Now().Add(-14*Day),
	))

	prdetails := try(gh.ViewPRsDetails(prs))

	slices.SortFunc(prdetails, func(a, b gh.PRDetails) int {
		return -a.UpdatedAt.Compare(b.UpdatedAt)
	})

	fmt.Println(prdetails.Style())
}

func try[T any](value T, err error) T {
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return value
}
