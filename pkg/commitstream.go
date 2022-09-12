package commitstream

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	DomainsFile         string
	DomainsList         map[string]bool
}

type Commit struct {
	Name  string
	Email struct {
		User   string
		Domain string
	}
	Repo    string
	Message string
}

func (cs *CommitStream) Start(handler Handler) {

	if cs.Filter.Email == "" && cs.Filter.Name == "" {
		cs.Filter.Enabled = false
	} else {
		cs.Filter.Enabled = true
	}
	if cs.Filter.DomainsFile != "" {
		cs.Filter.DomainsList = make(map[string]bool)
		f, err := os.Open(cs.Filter.DomainsFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			cs.Filter.DomainsList[scanner.Text()] = true
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		cs.Filter.Enabled = true

	}
	gh := GithubHandler{
		Cstream: cs,
	}

	//var handledCounter uint64

	var commitsChan = make(chan []Commit, 200)

	go func() {
		for range time.Tick(time.Second * 1) {
			if cs.Debug == true {
				s := cs.Stats

				msg := fmt.Sprintf("incoming: %d, processed: %d, accepted: %d, total: %d, chan sz:%d\n",
					s.IncomingRate, s.ProcessedRate, s.FilteredRate, s.Total, len(commitsChan))
				os.Stderr.WriteString(msg)
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

	if len(cs.Filter.DomainsList) != 0 {
		if ok := cs.Filter.DomainsList[c.Email.Domain]; ok {
			return true
		}
	}
	if cs.Filter.Email != "" {

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
