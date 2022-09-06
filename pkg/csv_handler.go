package commitstream

import (
	"encoding/csv"
	"os"
)

type CsvHander struct{}

type NoHandler struct{}

func (n NoHandler) Callback(commits []Commit) {
	//fmt.Println(c.Repo)
	//time.Sleep(time.Duration(rand.Intn(10)))
	return
}

func (h CsvHander) Callback(commits []Commit) {
	w := csv.NewWriter(os.Stdout)
	for _, c := range commits {
		email := c.Email.User + "@" + c.Email.Domain
		cOut := []string{c.Name, email, "https://github.com/" + c.Repo, c.Message}

		w.Write(cOut)
	}

	w.Flush()
}
