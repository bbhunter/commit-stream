package handlers

import (
	"encoding/csv"
	"os"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type CsvHander struct{}

type NoHandler struct{}

func (n NoHandler) Callback(commits []commit.CommitEvent) {
	//fmt.Println(c.Repo)
	//time.Sleep(time.Duration(rand.Intn(10)))
	return
}

func (h CsvHander) Callback(commits []commit.CommitEvent) {
	w := csv.NewWriter(os.Stdout)
	for _, c := range commits {
		email := c.AuthorEmail.User + "@" + c.AuthorEmail.Domain
		cOut := []string{c.UserName, email, "https://github.com/" + c.RepoName, c.Message}

		w.Write(cOut)
	}

	w.Flush()
}
