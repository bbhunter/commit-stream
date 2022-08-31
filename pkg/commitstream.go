package commitstream

import (
	"strings"
	"sync"
)

type CommitStream struct {
	mu            sync.Mutex
	GithubOptions *GithubOptions
	Filter        *Filter
}

type GithubOptions struct {
	AuthToken string
	Rate      int
}

type Filter struct {
	Email               string
	Name                string
	Enabled             bool
	IgnorePrivateEmails bool
	IncludeMessages     bool
	SearchAllCommits    bool
}

type Commit struct {
	Name    string
	Email   string
	Repo    string
	Message string
}

func (cs *CommitStream) Start(callback func(Commit)) {
	gh := GithubHandler{
		Cstream: cs,
	}

	var results = make(chan FeedResult)

	go func() {
		for result := range results {
			for e, n := range result.CommitAuthors {
				c := Commit{Name: n, Email: e, Repo: result.RepoName}
				if cs.Filter.IncludeMessages != false {
					c.Message = result.Message
				}

				if cs.isMatch(c) {
					cs.outputMatch(c, callback)
				}
			}
		}
	}()

	gh.Run(results)

}

func (cs *CommitStream) isMatch(c Commit) bool {

	if cs.Filter.IgnorePrivateEmails == true {
		if strings.Contains(c.Email, "@users.noreply.github.com") {
			return false
		}
	}

	if cs.Filter.Enabled == false {
		return true
	}

	result := false

	if cs.Filter.Email != "" {
		//fmt.Printf("checking %s against %s\n", email, fo.email)
		for _, e := range strings.Split(cs.Filter.Email, ",") {
			if strings.Contains(c.Email, strings.TrimSpace(e)) {
				result = true
			}
		}
	}

	if cs.Filter.Name != "" {
		for _, n := range strings.Split(cs.Filter.Name, ",") {
			if strings.Contains(c.Name, strings.TrimSpace(n)) {
				result = true
			}
		}
	}

	return result
}

func (cs *CommitStream) outputMatch(c Commit, callback func(Commit)) {
	//s := []string{c.name, c.email, c.repo}
	//tm := time.Now().UTC().Format("2006-01-02T15:04:05")

	cs.mu.Lock()
	callback(c)
	cs.mu.Unlock()
}
