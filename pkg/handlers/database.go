package handlers

import (
	"github.com/x1sec/commit-stream/pkg/commit"
	"github.com/x1sec/commit-stream/pkg/database"
)

type DatabaseHandler struct {
	Db database.Database
}

func (dh *DatabaseHandler) Callback(commits []commit.CommitEvent) {
	for _, c := range commits {
		email := c.AuthorEmail.User + "@" + c.AuthorEmail.Domain
		d := database.CommitEntry{
			AuthorName:  c.AuthorName,
			RepoName:    c.RepoName,
			UserName:    c.UserName,
			AuthorEmail: email,
			Message:     c.Message,
			SHA:         c.SHA,
		}
		dh.Db.Insert(d)
	}
}
