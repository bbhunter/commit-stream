package csv

import (
	"encoding/csv"
	"io"

	"log"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type CsvHander struct {
	out io.Writer
}

type NoHandler struct{}

func (n NoHandler) Callback(commits []commit.CommitEvent) {
	//fmt.Println(c.Repo)
	//time.Sleep(time.Duration(rand.Intn(10)))
	return
}

func NewCsvHandler(out io.Writer) CsvHander {
	log.Println("Using CSV handler")
	return CsvHander{out: out}
}

func (h CsvHander) Callback(commits []commit.CommitEvent) {
	w := csv.NewWriter(h.out)
	for _, c := range commits {
		email := c.AuthorEmail.User + "@" + c.AuthorEmail.Domain
		cOut := []string{c.AuthorName, email, "https://github.com/" + c.UserName + "/" + c.RepoName, c.Message}

		w.Write(cOut)
	}

	w.Flush()
}
