package gh

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

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
	IsDraft   bool      `json:"isDraft"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

type PullRequests []PullRequest

type ReviewDecision string

const (
	ReviewRequired   ReviewDecision = "REVIEW_REQUIRED"
	Accepted         ReviewDecision = "ACCEPTED"
	ChangesRequested ReviewDecision = "CHANGES_REQUESTED"
)

type MergeableStatus string

const (
	Mergeable   MergeableStatus = "MERGEABLE"
	Conflicting MergeableStatus = "CONFLICTING"
	Unknown     MergeableStatus = "UNKNOWN"
)

type PRDetails struct {
	PullRequest    `json:"-"`
	ReviewDecision ReviewDecision  `json:"reviewDecision"`
	Mergeable      MergeableStatus `json:"mergeable"`
}

func (pr PRDetails) Style() string {
	color := "#ffffff"
	icon := ""
	switch {
	case pr.IsDraft:
		color = "#777777"
		icon = ""
	case pr.State == Closed:
		color = "#C53211"
		icon = ""
	case pr.State == Merged:
		color = "#853CEA"
		icon = ""
	case pr.State == Open:
		color = "#0FBF3E"
		icon = ""
	}

	green := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(termenv.ANSIBrightGreen))
	red := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(termenv.ANSIBrightRed))

	reviewDecision := ""
	switch pr.ReviewDecision {
	case Accepted:
		reviewDecision = green.Render(" ")
	case ChangesRequested:
		reviewDecision = red.Render(" ")
	case ReviewRequired:
		reviewDecision = ""
	}

	mergeableStatus := ""
	switch pr.Mergeable {
	case Mergeable:
		mergeableStatus = green.Render(" ")
	case Conflicting:
		mergeableStatus = red.Render(" ")
	case Unknown:
		mergeableStatus = ""
	}

	card := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		PaddingLeft(2).
		PaddingRight(2).
		MarginLeft(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		Width(80).
		AlignHorizontal(lipgloss.Left)

	neutral := lipgloss.NewStyle().Foreground(lipgloss.Color("#bbb"))

	return card.Render(fmt.Sprintf("%s%s%s %s \n%s | %s",
		icon,
		reviewDecision,
		mergeableStatus,
		pr.Title,
		neutral.Render(pr.Author.Login),
		neutral.Render(pr.URL))) + "\n"
}

type PRDetailsList []PRDetails

func (prs PRDetailsList) Style() string {
	bob := strings.Builder{}
	for i := range prs {
		bob.WriteString(prs[i].Style())
	}
	return bob.String()
}
