package github

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/go-github/github"
	ggithub "github.com/google/go-github/github"
	"github.com/x1sec/commit-stream/pkg/commit"
	"github.com/x1sec/commit-stream/pkg/stats"
	"golang.org/x/oauth2"
)

type GithubOptions struct {
	AuthToken string
	Rate      int
}
type GithubHandler struct {
	session Session
	Options *GithubOptions
}

type Session struct {
	client *ggithub.Client
	ctx    context.Context
}

type FeedResult struct {
	CommitAuthors map[string]string
	RepoName      string
	RepoURL       string
	Message       string
}

func (gh *GithubHandler) checkResponseError(err error, resp *github.Response) bool {
	if _, ok := err.(*ggithub.RateLimitError); ok {
		log.Println("Hit rate limit. Reset: %s\n", resp.Rate.Reset)
		time.Sleep(time.Until(resp.Rate.Reset.Time))
		return true
	}
	if _, ok := err.(*ggithub.AbuseRateLimitError); ok {
		fmt.Fprintf(os.Stderr, "Abuse detected!\n")
		os.Exit(1)
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		fmt.Fprintf(os.Stderr, "Timeout occured, sleeping for 5 seconds...\n")
		time.Sleep(5 * time.Second)
		return true
	}

	if err, r := err.(*ggithub.ErrorResponse); r {
		switch statusCode := err.Response.StatusCode; statusCode {
		case 401:
			fmt.Fprintf(os.Stderr, "401 - Error with authentication token provided.\n")
			os.Exit(1)
		case 502:
			// Handle 502 sleeping for file seconds before retrying
			fmt.Fprintf(os.Stderr, "502 - Bad Gateway, sleeping for 5 seconds... \n")
			time.Sleep(5 * time.Second)
			return true
		default:
			fmt.Fprintf(os.Stderr, err.Error())
			return true
		}
		return false

	}

	return false
}

func (gh *GithubHandler) Run(results chan<- []commit.CommitEvent, stats *stats.ProcessingStats, searchAllCommits bool) {
	options := gh.Options

	gh.session.ctx = context.Background()
	lc, cancel := context.WithCancel(gh.session.ctx)

	defer cancel()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: options.AuthToken},
	)
	tc := oauth2.NewClient(gh.session.ctx, ts)

	gh.session.client = ggithub.NewClient(tc)
	for {
		opt := &ggithub.ListOptions{PerPage: 300}
		for {

			events, resp, err := gh.session.client.Activity.ListEvents(lc, opt)

			if gh.checkResponseError(err, resp) {
				continue
			}

			var commits []commit.CommitEvent
			for _, e := range events {

				if *e.Type == "PushEvent" {

					p, _ := e.ParsePayload()

					q := p.(*ggithub.PushEvent)

					for _, r := range q.Commits {
						var commit commit.CommitEvent
						commit.Timestamp = time.Now()
						userRepo := *e.GetRepo().Name
						_userRepo := strings.Split(userRepo, "/")
						commit.UserName = _userRepo[0]
						commit.RepoName = _userRepo[1]
						commit.SHA = *r.SHA
						commit.Message = *r.Message
						commit.AuthorName = *r.GetAuthor().Name

						email := *r.GetAuthor().Email

						if strings.Contains(email, "@") {
							parts := strings.Split(email, "@")
							commit.AuthorEmail.User = parts[0]
							commit.AuthorEmail.Domain = parts[1]
						} else {
							commit.AuthorEmail.User = email
						}

						atomic.AddUint32(&stats.IncomingRate, 1)

						commits = append(commits, commit)
						if searchAllCommits == false {
							break
						}
					}

				}
			}
			results <- commits

			//fmt.Fprintf(os.Stderr, "\r%d/%d remaining\n", resp.Rate.Remaining, resp.Rate.Limit)
			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage

			time.Sleep(time.Second * 1)

		}

		time.Sleep(time.Second * time.Duration(1))

	}
}
