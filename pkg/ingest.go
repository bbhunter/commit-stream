package commitstream

import (
	"fmt"
	"os"
	"strings"

	"time"
	"sync"
)

type FilterOptions struct {
	Email string
	Name string
	Enabled bool

} 

type commit struct {
	name string
	email string
	repo string
}

type FlagOption byte

const (
	DATE FlagOption = 1 << iota
	NAME 
	EMAIL
	REPO
)

var mu sync.Mutex 


func DoIngest(streamOpt StreamOptions, fo FilterOptions, fmtFlags FlagOption, callback func([]string)) {


	var results = make(chan FeedResult)

	go func() {
		for result := range results {
			for e, n := range result.CommitAuthors {
				c := commit{n, e, result.RepoURL}
				if isMatch(c,fo) {
					outputMatch(c, fmtFlags, callback)
				}			
			}
		}
	}()
	
	//streamOpt := commitstream.StreamOptions{AuthToken : authToken, SearchAllCommits: searchAllCommits}
	Run(streamOpt, results)

}

func ParseOutputOption(f *FlagOption, format string) {

	for _, c := range format {
		switch string(c) {
			case "d" :
				*f |= DATE
			case "n" :
				*f |= NAME
			case "e" :
				*f |= EMAIL
			case "r" :
				*f |= REPO
			default :
				fmt.Fprintf(os.Stderr, "Invalid output modifier specified: %s\n", string(c))
				os.Exit(1)
		}
	}
}

func isMatch(c commit, fo FilterOptions) bool {
	
	if fo.Enabled == false {
		return true
	} 

	result := false
	

	if fo.Email != "" {
		//fmt.Printf("checking %s against %s\n", email, fo.email)
		for _, e := range strings.Split(fo.Email, ",") {
			if strings.Contains(c.email, strings.TrimSpace(e)) {
				result = true
			}
		}
	}

	if fo.Name != "" {
		for _, n := range strings.Split(fo.Name, ",") {
			if strings.Contains(c.name, strings.TrimSpace(n)) {
				result = true
			}
		}
	}					

	return result
}

func outputMatch(c commit, f FlagOption, callback func([]string)) {
	var s []string

	if f == 0 {
		f = DATE | NAME | EMAIL | REPO
	}

	if f & DATE != 0 {
		tm := time.Now().UTC().Format("2006-01-02T15:04:05")
		s = append(s, tm)
	}

	if f & NAME != 0 {
		s = append(s, c.name)
	}
	if f & EMAIL != 0 {
		s = append(s, c.email)
	}
	if f & REPO != 0 {
		s = append(s, c.repo)
	}

	mu.Lock()
	callback(s)
	mu.Unlock()
}

