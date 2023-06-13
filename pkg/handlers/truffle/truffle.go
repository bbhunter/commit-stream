package truffle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/slack-go/slack"
	"github.com/x1sec/commit-stream/pkg/commit"
	slackconf "github.com/x1sec/commit-stream/pkg/handlers/slack"
)

type Truffle struct {
	DetectorName   string      `json:"DetectorName"`
	DetectorType   int         `json:"DetectorType"`
	ExtraData      interface{} `json:"ExtraData"`
	Raw            string      `json:"Raw"`
	Redacted       string      `json:"Redacted"`
	SourceID       int         `json:"SourceID"`
	SourceMetadata struct {
		Data struct {
			Github struct {
				Commit     string `json:"commit"`
				Email      string `json:"email"`
				File       string `json:"file"`
				Line       int    `json:"line"`
				Link       string `json:"link"`
				Repository string `json:"repository"`
				Timestamp  string `json:"timestamp"`
			} `json:"Github"`
		} `json:"Data"`
	} `json:"SourceMetadata"`
	SourceName     string      `json:"SourceName"`
	SourceType     int         `json:"SourceType"`
	StructuredData interface{} `json:"StructuredData"`
	Verified       bool        `json:"Verified"`
}

type TruffleHandler struct {
	Path           string
	MaxWorkers     int
	DroppedCommits uint64
	LogFile        string
	GithubToken    string
	lastList       map[string]time.Time
	queue          chan commit.CommitEvent
	SlackConf      slackconf.SlackHandler
}

func NewTruffleHandler(path string, maxWorkers int, githubToken string) *TruffleHandler {
	log.Println("Using truffle handler")

	ts := &TruffleHandler{
		Path:        path,
		MaxWorkers:  maxWorkers,
		GithubToken: githubToken,
	}
	ts.StartWorkers()
	return ts
}

func (h *TruffleHandler) parseJson(entry string) (Truffle, error) {
	var te Truffle
	err := json.Unmarshal([]byte(entry), &te)
	if err != nil {
		return te, err
	}
	return te, nil
}
func (h *TruffleHandler) Run(worker int, commit commit.CommitEvent) {
	for repo, tm := range h.lastList {
		if time.Now().After(tm) {
			log.Printf("[%s] expired, removing", commit.RepoName)
			delete(h.lastList, repo)
		}
	}
	k := commit.UserName + commit.RepoName
	if _, ok := h.lastList[k]; ok {
		return
	}
	h.lastList[k] = time.Now().Add(time.Minute)

	if _, err := os.Stat(h.Path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	url := "https://github.com/" + commit.AuthorName + "/" + commit.RepoName
	log.Printf("[%d] Running truffle for: %s", worker, url)
	cmd := exec.Command(h.Path, "github", "--repo="+url, "--token="+h.GithubToken, "--json")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("error in cmd.Start()")
		log.Panicln(err)
		return
	}
	body, err := ioutil.ReadAll(stdout)

	if err != nil {
		fmt.Println("error in ioutil.ReadAll()")
		log.Println(err)
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Error in cmd.Wait()")
		body, _ := ioutil.ReadAll(stderr)
		fmt.Println(body)
		log.Println(err)
		return
	}
	if len(body) > 1 {
		// TODO: Normal output
		log.Println(string(body))
		lines := strings.Split(string(body), "\n")
		log.Printf("Parsing %d truffles... ", len(lines))
		for _, line := range lines {
			truffle, err := h.parseJson(line)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Posting truffle to Slack: " + truffle.SourceMetadata.Data.Github.Repository)
				h.SendSlack(truffle)
			}
		}

	}
	log.Printf("[%d] done", worker)

}

func (h *TruffleHandler) StartWorkers() {
	h.queue = make(chan commit.CommitEvent, 1000)
	h.lastList = make(map[string]time.Time)

	// spin up workers
	for i := 0; i < h.MaxWorkers; i++ {
		go func(j int) {
			for c := range h.queue {
				h.Run(j, c)
			}
		}(i)
	}

	// queue stats worder
	go func() {
		for {
			if h.DroppedCommits > 0 {
				fmt.Printf("%d\n", h.DroppedCommits)
			}
			//fmt.Printf("queue size: %d\n", len(queue))
			time.Sleep(time.Second * 5)
		}
	}()
}

func (h *TruffleHandler) Callback(commits []commit.CommitEvent) {
	for _, c := range commits {
		select {
		case h.queue <- c:
		default:
			//fmt.Println("Queue exhausted")
			atomic.AddUint64(&h.DroppedCommits, 1)
		}

	}
	//close(queue)
}

func (s TruffleHandler) SendSlack(t Truffle) {

	secret := t.Raw
	if len(secret) > 30 {
		secret = secret[0:30] + ".. <cut> .."
	}
	client := slack.New(s.SlackConf.Token, slack.OptionDebug(false))
	attachment := slack.Attachment{
		Pretext: "Commit-stream message",
		Color:   "#36a64f",
		Text:    t.DetectorName,
		Fields: []slack.AttachmentField{
			{
				Title: "Repository",
				Value: t.SourceMetadata.Data.Github.Email,
			},
			{
				Title: "File",
				Value: t.SourceMetadata.Data.Github.Link,
			},
			{
				Title: "Line",
				Value: strconv.Itoa(t.SourceMetadata.Data.Github.Line),
			},
			{
				Title: "Secret",
				Value: string(secret),
			},
		},
	}
	_, timestamp, err := client.PostMessage(
		s.SlackConf.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		log.Println(err)
	}
	log.Printf("Slack message sent at %s", timestamp)
}
