/*
commit-stream
Author: https://twitter.com/haxrob

See LICENSE
*/

package commitstream

import (
	"strings"
	"sync"
)

type FilterOptions struct {
	Email               string
	Name                string
	Enabled             bool
	IgnorePrivateEmails bool
}

type Commit struct {
	Name  string
	Email string
	Repo  string
}

type Callback interface {
	Run(c Commit)
}

var mu sync.Mutex

func DoIngest(streamOpt StreamOptions, fo FilterOptions, callback func(Commit)) {

	var results = make(chan FeedResult)

	go func() {
		for result := range results {
			for e, n := range result.CommitAuthors {
				c := Commit{n, e, result.RepoURL}
				if isMatch(c, fo) {
					outputMatch(c, callback)
				}
			}
		}
	}()

	Run(streamOpt, results)

}

func isMatch(c Commit, fo FilterOptions) bool {

	if fo.IgnorePrivateEmails == true {
		if strings.Contains(c.Email, "@users.noreply.github.com") {
			return false
		}
	}

	if fo.Enabled == false {
		return true
	}

	result := false

	if fo.Email != "" {
		//fmt.Printf("checking %s against %s\n", email, fo.email)
		for _, e := range strings.Split(fo.Email, ",") {
			if strings.Contains(c.Email, strings.TrimSpace(e)) {
				result = true
			}
		}
	}

	if fo.Name != "" {
		for _, n := range strings.Split(fo.Name, ",") {
			if strings.Contains(c.Name, strings.TrimSpace(n)) {
				result = true
			}
		}
	}

	return result
}

func outputMatch(c Commit, callback func(Commit)) {
	//s := []string{c.name, c.email, c.repo}
	//tm := time.Now().UTC().Format("2006-01-02T15:04:05")

	mu.Lock()
	callback(c)
	mu.Unlock()
}
