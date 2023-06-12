package truffle

import (
	"github.com/x1sec/commit-stream/pkg/handlers/slack"
)

func NewTruffleSlackHander(slackToken string, slackChannel string, trufflePath string, truffleWorkers int, githubToken string) *TruffleHandler {
	s := slack.NewSlackHandler(slackToken, slackChannel)
	h := TruffleHandler{
		Slack:       s,
		Path:        trufflePath,
		MaxWorkers:  truffleWorkers,
		GithubToken: githubToken,
	}
	h.StartWorkers()
	return &h
}
