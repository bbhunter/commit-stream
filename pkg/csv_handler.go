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
		cOut := []string{c.Name, c.Email, "https://github.com/" + c.Repo, c.Message}

		w.Write(cOut)
	}

	w.Flush()
}
