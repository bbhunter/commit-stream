package commitstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type ElasticHandler struct {
	RemoteURI string
	Username  string
	Password  string
}

func (e ElasticHandler) Callback(commits []Commit) {
	e.ImportBulk(commits)
}

func (e *ElasticHandler) Import(commit Commit) {

	path := "/api/commits/_doc"

	data, err := json.Marshal(commit)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", e.RemoteURI+path, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	if e.Username != "" {
		req.SetBasicAuth(e.Username, e.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)

		log.Println(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	} else {
		resp.Body.Close()
	}
}

func (e *ElasticHandler) ImportBulk(commits []Commit) {

	path := "/api/_bulk"

	var entry string
	for _, commit := range commits {
		data, err := json.Marshal(commit)
		if err == nil {
			entry = entry + `{ "index" : { "_index" : "commits" } }`
			entry = entry + "\n" + string(data) + "\n"
		} else {
			log.Fatal(err)
		}
	}

	req, err := http.NewRequest("POST", e.RemoteURI+path, strings.NewReader(entry))
	if err != nil {
		log.Fatal(err)
	}
	if e.Username != "" {
		req.SetBasicAuth(e.Username, e.Password)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		defer resp.Body.Close()
		log.Println(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}
}
