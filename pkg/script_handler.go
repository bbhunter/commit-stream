package commitstream

import (
	"fmt"
	"log"
	"os/exec"
	"sync/atomic"
	"time"
)

type ScriptHandler struct {
	Path           string
	MaxWorkers     int
	DroppedCommits uint64
}

func (h ScriptHandler) Run(worker int, commit Commit) {
	cmd := exec.Command(h.Path, commit.Repo)
	//fmt.Printf("[%d] start\n", worker)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("[%d] end\n", worker)

}

func (h ScriptHandler) Callback(commits []Commit) {

	queue := make(chan Commit, h.MaxWorkers)

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
