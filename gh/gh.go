package gh

import (
	"encoding/json"
	"fmt"
	"time"

	cli "github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/rs/zerolog/log"
)

var restClient, _ = api.DefaultRESTClient()

func GetUser() (string, error) {
	response := struct{ Login string }{}
	err := restClient.Get("user", &response)
	if err != nil {
		return "", err
	}

	return response.Login, nil
}

func SearchPRs(authors []string, createdAfter time.Time) (PullRequests, error) {
	args := makeSearchPRArgs(authors, createdAfter)

	log.Info().Msg("searching PRs")
	stdout, _, err := cli.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("cli: %w", err)
	}

	var prs PullRequests
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func makeSearchPRArgs(authors []string, createdAfter time.Time) []string {
	args := []string{
		"search", "prs",
		"--created", ">=" + createdAfter.Format(time.DateOnly),
		"--json", "author,createdAt,updatedAt,title,state,url",
		"--",
	}

	// --- filter authors ---
	for i, author := range authors {
		if i > 0 {
			args = append(args, "OR")
		}
		args = append(args, "author:"+author)
	}

	return args
}
