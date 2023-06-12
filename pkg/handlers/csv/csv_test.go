package csv_test

import (
	"bytes"
	"testing"

	"github.com/x1sec/commit-stream/pkg/commit"
	"github.com/x1sec/commit-stream/pkg/handlers/csv"
)

func TestCallback(t *testing.T) {
	var mockOut bytes.Buffer
	h := csv.NewCsvHandler(&mockOut)
	c := commit.CommitEvent{
		RepoName:   "reponame",
		UserName:   "username",
		AuthorName: "authorname",
		Message:    "message",
		SHA:        "123456789",
		EventType:  "push",
	}
	c.AuthorEmail.Domain = "test.com"
	c.AuthorEmail.User = "tester"
	ca := []commit.CommitEvent{c}
	h.Callback(ca)
	s := mockOut.String()
	expected := "authorname,tester@test.com,https://github.com/username/reponame,message\n"
	//t.Logf("expected:\n%s\n", expected)
	//t.Logf("got:\n%s\n", s)
	if s != expected {
		t.Error("wrong csv output")
	}

}
