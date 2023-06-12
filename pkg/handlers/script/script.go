package script

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type ScriptHandler struct {
	Path           string
	MaxWorkers     int
	DroppedCommits uint64
	LogFile        string
}

func NewScriptHandler(path string, maxWorkers int, logFile string) *ScriptHandler {
	return &ScriptHandler{
		Path:       path,
		MaxWorkers: maxWorkers,
		LogFile:    logFile,
	}
}

func (h ScriptHandler) Run(worker int, commit commit.CommitEvent) {
	if _, err := os.Stat(h.Path); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	//s := strings.Split(commit.Repo, "/")

	cmd := exec.Command(h.Path, commit.UserName, commit.RepoName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("a")
		log.Panicln(err)
		return
	}
	body, err := ioutil.ReadAll(stdout)
	//fmt.Println(body)
	if err != nil {
		fmt.Println("d")
		log.Println(err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("b")
		body, _ := ioutil.ReadAll(stderr)
		fmt.Println(body)
		log.Println(err)
		return
	}

	f, err := os.OpenFile(h.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	fmt.Println(h.LogFile)
	if err != nil {
		fmt.Println("c")
		log.Println(err)
		return
	}
	log.Println(string(body))
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		fmt.Println("e")
		log.Println(err)
	}

}

func (h ScriptHandler) Callback(commits []commit.CommitEvent) {

	queue := make(chan commit.CommitEvent, h.MaxWorkers)

	for i := 0; i < h.MaxWorkers; i++ {
		go func(j int) {
			for c := range queue {
				h.Run(j, c)
			}
		}(i)
	}
	go func() {
		for {
			if h.DroppedCommits > 0 {
				fmt.Printf("%d\n", h.DroppedCommits)
			}
			//fmt.Printf("queue size: %d\n", len(queue))
			time.Sleep(time.Second * 5)
		}
	}()
	for _, c := range commits {
		select {
		case queue <- c:
		default:
			//fmt.Println("Queue exhausted")
			atomic.AddUint64(&h.DroppedCommits, 1)
		}

	}
	close(queue)
}
