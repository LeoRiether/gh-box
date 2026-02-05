package gh

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LeoRiether/gh-box/workers"
	cli "github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
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

type SearchOptions struct {
	Authors      []string
	Organization string
	CreatedAfter *time.Time
	UpdatedAfter *time.Time
	State        PRState
}

func SearchPRs(opts SearchOptions) (PullRequests, error) {
	args := makeSearchPRArgs(opts)

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

func makeSearchPRArgs(opts SearchOptions) []string {
	args := []string{
		"search", "prs",
		"--json", "author,createdAt,updatedAt,title,state,isDraft,url",
	}

	// --- filter based on the date ---
	if opts.CreatedAfter != nil {
		args = append(args, "--created", ">="+opts.CreatedAfter.Format(time.DateOnly))
	}
	if opts.UpdatedAfter != nil {
		args = append(args, "--updated", ">="+opts.UpdatedAfter.Format(time.DateOnly))
	}

	args = append(args, "--")

	// --- filter authors ---
	if len(opts.Authors) > 0 {
		for i, author := range opts.Authors {
			if i > 0 {
				args = append(args, "OR")
			}
			args = append(args, "author:"+author)

			if i == 0 {
				args[len(args)-1] = "(" + args[len(args)-1]
			}
			if i == len(opts.Authors)-1 {
				args[len(args)-1] = args[len(args)-1] + ")"
			}
		}
	}

	// --- filter organization ---
	if opts.Organization != "" {
		args = append(args, "org:"+opts.Organization)
	}

	switch opts.State {
	case Open:
		args = append(args, "state:open")
	case Closed:
		args = append(args, "state:closed")
	case Merged:
		args = append(args, "is:merged")
	}

	return args
}

func ViewPRsDetails(prs PullRequests) (PRDetailsList, error) {
	pool := workers.NewPool(16, func(pr PullRequest) (PRDetails, error) {
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
