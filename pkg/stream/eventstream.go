package stream

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/x1sec/commit-stream/pkg/commit"
	"github.com/x1sec/commit-stream/pkg/filter"
	"github.com/x1sec/commit-stream/pkg/github"
	"github.com/x1sec/commit-stream/pkg/handlers"
	"github.com/x1sec/commit-stream/pkg/stats"
)

type CommitEventStream struct {
	mu            sync.Mutex
	GithubOptions *github.GithubOptions
	Filter        *filter.Filter
	Stats         stats.ProcessingStats
	Debug         bool
}

func (cs *CommitEventStream) Start(handler handlers.Handler) {

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
		log.Printf("Loaded %d domains from %s\n", len(cs.Filter.DomainsList), cs.Filter.DomainsFile)
		cs.Filter.Enabled = true

	}
	gh := github.GithubHandler{
		Options: cs.GithubOptions,
	}

	//var handledCounter uint64

	var commitsChan = make(chan []commit.CommitEvent, 10000)

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
			var filteredCommitEvents []commit.CommitEvent
			for _, commit := range commits {
				if cs.Filter.IncludeMessages == false {
					commit.Message = ""
				}
				atomic.AddUint32(&cs.Stats.ProcessedRate, 1)
				if filter.Filtered(commit, *cs.Filter) {
					atomic.AddUint32(&cs.Stats.FilteredRate, 1)
					filteredCommitEvents = append(filteredCommitEvents, commit)
				}
			}
			cs.execHandler(filteredCommitEvents, handler)
		}

	}()

	gh.Run(commitsChan, &cs.Stats, cs.Filter.SearchAllCommitEvents)

}

func (cs *CommitEventStream) execHandler(commits []commit.CommitEvent, handler handlers.Handler) {
	cs.mu.Lock()
	handler.Callback(commits)
	cs.mu.Unlock()
}
