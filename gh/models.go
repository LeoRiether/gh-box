package gh

import "time"

type Author struct {
	Login string `json:"login"`
}

type PRState string

const (
	Open   PRState = "open"
	Merged PRState = "merged"
	Closed PRState = "closed"
)

type PullRequest struct {
	Author    Author    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	State     PRState   `json:"state"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

type PullRequests []PullRequest
