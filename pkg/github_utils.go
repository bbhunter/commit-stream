package commitstream

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubUtil struct {
	session Session
	Token   string
}

func (g *GithubUtil) newSession() {
	g.session.ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)
	tc := oauth2.NewClient(g.session.ctx, ts)
	g.session.client = github.NewClient(tc)
}

func (g *GithubUtil) GetEmailsByRepo(user string, repo string) (emails []string, err error) {
	unique := make(map[string]bool)
	if g.session.client == nil {
		g.newSession()
	}
	commits, response, err := g.session.client.Repositories.ListCommits(context.Background(), user, repo, nil)
	if err != nil && response.StatusCode != 200 {
		return nil, err
	}
	for _, commit := range commits {
		email := *commit.Commit.Author.Email
		if strings.Contains(email, "@users.noreply.github.com") {
			continue
		}
		if _, value := unique[email]; !value {
			unique[email] = true
			emails = append(emails, email)
		}
	}
	return emails, nil
}

func (g *GithubUtil) GetLastCommitAuthor(user string, repo string) (commit Commit, err error) {
	if g.session.client == nil {
		g.newSession()
	}
	commits, response, err := g.session.client.Repositories.ListCommits(context.Background(), user, repo, nil)
	if err != nil && response.StatusCode != 200 {
		return Commit{}, err
	}
	if len(commits) == 0 {
		return Commit{}, nil
	}
	lastCommit := commits[0]
	email := *lastCommit.Commit.Author.Email
	name := *lastCommit.Commit.Author.Name
	message := *lastCommit.Commit.Message
	if strings.Contains(email, "@") {
		parts := strings.Split(email, "@")
		commit.Email.User = parts[0]
		commit.Email.Domain = parts[1]
	} else {
		commit.Email.User = email
	}
	commit.Name = name
	commit.Message = message
	return commit, nil
}
