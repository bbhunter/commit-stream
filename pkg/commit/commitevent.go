package commit

import (
	"time"
)

type CommitEvent struct {
	AuthorEmail struct {
		User   string
		Domain string
	}
	RepoName   string
	UserName   string
	AuthorName string
	Message    string
	Timestamp  time.Time
	SHA        string
	EventType  string
}
