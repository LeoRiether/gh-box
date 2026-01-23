package gh

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Author struct {
	Login string `json:"login"`
}

type PRState string

const (
	Draft  PRState = "draft"
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

func (pr PullRequest) Style() string {
	color := "#ffffff"
	icon := ""
	switch pr.State {
	case Closed:
		color = "#C53211"
		icon = ""
	case Merged:
		color = "#853CEA"
		icon = ""
	case Open:
		color = "#0FBF3E"
		icon = ""
	case Draft:
		color = "#777777"
		icon = ""
	}

	style := lipgloss.NewStyle().
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

	neutral := lipgloss.NewStyle().Foreground(lipgloss.Color("#999"))

	return style.Render(fmt.Sprintf("%s %s\n%s",
		icon,
		pr.Title,
		neutral.Render(pr.Author.Login))) + "\n"
}

type PullRequests []PullRequest

func (prs PullRequests) Style() string {
	bob := strings.Builder{}
	for i := range prs {
		bob.WriteString(prs[i].Style())
	}
	return bob.String()
}
