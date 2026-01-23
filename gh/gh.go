package gh

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LeoRiether/gh-box/workers"
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

func SearchPRs(authors []string, org string, createdAfter time.Time) (PullRequests, error) {
	args := makeSearchPRArgs(authors, org, createdAfter)

	log.Info().Msg("searching PRs")
	stdout, stderr, err := cli.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("cli: %w. stderr: %v", err, stderr.String())
	}

	var prs PullRequests
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, err
	}

	return prs, nil
}

func makeSearchPRArgs(authors []string, org string, createdAfter time.Time) []string {
	args := []string{
		"search", "prs",
		// "--created", ">=" + createdAfter.Format(time.DateOnly),
		"--json", "author,createdAt,updatedAt,title,state,isDraft,url",
		"--",
	}

	// --- filter authors ---
	for i, author := range authors {
		if i > 0 {
			args = append(args, "OR")
		}
		args = append(args, "author:"+author)

		if i == 0 {
			args[len(args)-1] = "(" + args[len(args)-1]
		} else if i == len(authors)-1 {
			args[len(args)-1] = args[len(args)-1] + ")"
		}
	}

	// --- filter organization ---
	if org != "" {
		args = append(args, "org:"+org)
	}

	return args
}

func ViewPRsDetails(prs PullRequests) (PRDetailsList, error) {
	log.Info().Int("prs", len(prs)).Msg("fetching PR details")

	pool := workers.NewPool(8, func(pr PullRequest) (PRDetails, error) {
		stdout, stderr, err := cli.Exec(
			"pr", "view", pr.URL,
			"--json", "reviewDecision,mergeable")
		if err != nil {
			return PRDetails{}, fmt.Errorf("cli: %w. stderr: %v", err, stderr.String())
		}

		var details PRDetails
		if err := json.Unmarshal(stdout.Bytes(), &details); err != nil {
			return PRDetails{}, err
		}

		details.PullRequest = pr
		return details, nil
	})

	return pool.Process(prs)
}
