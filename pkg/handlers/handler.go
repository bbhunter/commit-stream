package handlers

import "github.com/x1sec/commit-stream/pkg/commit"

type Handler interface {
	Callback([]commit.CommitEvent)
}
