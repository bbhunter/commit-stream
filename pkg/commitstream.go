package commitstream

import (
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Handler interface {
	Callback([]Commit)
}

type ProcessingStats struct {
	IncomingRate  uint32
	ProcessedRate uint32
	FilteredRate  uint32
	Total         uint32
}
type CommitStream struct {
	mu            sync.Mutex
	GithubOptions *GithubOptions
	Filter        *Filter
	Stats         ProcessingStats
	Debug         bool
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
	Name  string
	Email struct {
		User   string
		Domain string
	}
	Repo    string
	Message string
	SHA     string
}

func (cs *CommitStream) Start(handler Handler) {
	gh := GithubHandler{
		Cstream: cs,
	}

	//var handledCounter uint64

	var commitsChan = make(chan []Commit, 200)

	go func() {
		for range time.Tick(time.Second * 1) {
			if cs.Debug == true {
				s := cs.Stats
				log.Printf("incoming: %d, processed: %d, accepted: %d, total: %d, chan sz:%d\n",
					s.IncomingRate, s.ProcessedRate,
					s.FilteredRate, s.Total, len(commitsChan))
			}
			atomic.AddUint32(&cs.Stats.Total, cs.Stats.FilteredRate)
			atomic.StoreUint32(&cs.Stats.ProcessedRate, 0)
			atomic.StoreUint32(&cs.Stats.FilteredRate, 0)
			atomic.StoreUint32(&cs.Stats.IncomingRate, 0)
		}
	}()

	go func() {
		for commits := range commitsChan {
			var filteredCommits []Commit
			for _, commit := range commits {
				if cs.Filter.IncludeMessages == false {
					commit.Message = ""
				}
				atomic.AddUint32(&cs.Stats.ProcessedRate, 1)
				if cs.filter(commit) {
					atomic.AddUint32(&cs.Stats.FilteredRate, 1)
					filteredCommits = append(filteredCommits, commit)
				}
			}
			cs.execHandler(filteredCommits, handler)
		}

	}()

	gh.Run(commitsChan)

}

func (cs *CommitStream) filter(c Commit) bool {

	if cs.Filter.IgnorePrivateEmails == true {
		if strings.Contains(c.Email.Domain, "users.noreply.github.com") {
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
			email := c.Email.User + "@" + c.Email.Domain
			if strings.Contains(email, strings.TrimSpace(e)) {
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

func (cs *CommitStream) execHandler(commits []Commit, handler Handler) {

	//tm := time.Now().UTC().Format("2006-01-02T15:04:05")

	cs.mu.Lock()
	handler.Callback(commits)
	cs.mu.Unlock()
}
