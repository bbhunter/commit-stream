package slack

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/x1sec/commit-stream/pkg/commit"
)

type SlackHandler struct {
	Token     string
	ChannelID string
}

func NewSlackHandler(token string, channelID string) SlackHandler {
	log.Println("Using slack handler")
	new := SlackHandler{
		Token:     token,
		ChannelID: channelID,
	}
	return new
}

func (s SlackHandler) Callback(commits []commit.CommitEvent) {
	for _, c := range commits {
		s.PostMessage(c)
	}
}

func (s SlackHandler) PostMessage(commit commit.CommitEvent) {
	client := slack.New(s.Token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext:    "commit-stream: incoming commit",
		Color:      "#36a64f",
		Text:       commit.AuthorEmail.Domain,
		AuthorName: commit.UserName + " / " + commit.RepoName,
		AuthorLink: "https://github.com/" + commit.UserName + "/" + commit.RepoName,
		Fields: []slack.AttachmentField{
			{
				Title: "Author Name",
				Value: commit.AuthorName,
				Short: true,
			},
			{
				Title: "Author Email",
				Value: commit.AuthorEmail.User + "@" + commit.AuthorEmail.Domain,
				Short: true,
			},
		},
	}
	_, _, err := client.PostMessage(
		s.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		log.Fatal("Failure with slack handler: " + err.Error())
	}
}
