package jira

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

const baseURL = "https://reload.atlassian.net"

// GetIssues gets issues.
func GetIssues() ([]jira.Issue, error) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("TRIAGEBOT_JIRA_USER"),
		Password: os.Getenv("TRIAGEBOT_JIRA_PASS"),
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, err
	}

	jql := fmt.Sprintf("filter = %s", os.Getenv("TRIAGEBOT_JIRA_FILTER"))

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

// FormatIssues formats issues as Markdown.
func FormatIssues(issues []jira.Issue) string {
	output := make([]string, 0, len(issues))

	for _, issue := range issues {
		// We HTML and URL encode the dash in the issue key to
		// hide the issue from Hubot (otherwise Hubot will
		// spam with follow up comments).
		issueKeyHTML := strings.Replace(issue.Key, "-", "&#x2D;", 1)
		issueKeyURL := strings.Replace(issue.Key, "-", "%2d", 1)
		output = append(
			output,
			fmt.Sprintf("* [%s](%s/browse/%s) - %s", issueKeyHTML, baseURL, issueKeyURL, issue.Fields.Summary),
		)
	}

	return strings.Join(output, "\n")
}
