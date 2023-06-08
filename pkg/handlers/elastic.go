package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type ElasticHandler struct {
	RemoteURI    string
	Username     string
	Password     string
	NoDuplicates bool
	UseZincAwsS3 bool
}

func (e ElasticHandler) Callback(commits []commit.CommitEvent) {
	//e.Import((commits[0]))
	e.ImportBulk(commits)
}

func (e *ElasticHandler) DoPost(path string, data string) {
	req, err := http.NewRequest("POST", e.RemoteURI+path, strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	if e.Username != "" {
		req.SetBasicAuth(e.Username, e.Password)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(resp.StatusCode)
			log.Println(err)
		}
		fmt.Println(string(body))
	}
	defer resp.Body.Close()
}

func (e *ElasticHandler) Setup() {
	if e.UseZincAwsS3 == true {
		log.Println("Using AWS S3 storage")
		path := "/api/index"
		storageRequest := `{"name": "commits", "storage_type": "s3"}`
		e.DoPost(path, storageRequest)
	} else {
		fmt.Println("Using local storage")
	}

}
func (e *ElasticHandler) Import(commit commit.CommitEvent) {

	path := "/api/commits/_doc"

	data, err := json.Marshal(commit)
	if err != nil {
		log.Fatal(err)
	}

	e.DoPost(path, string(data))
}

func (e *ElasticHandler) ImportBulk(commits []commit.CommitEvent) {

	path := "/api/_bulk"

	var entry string

	for _, commit := range commits {
		data, err := json.Marshal(commit)
		if err == nil {
			if e.NoDuplicates == true {
				key := []byte(commit.RepoName + commit.AuthorEmail.Domain)
				id := fmt.Sprintf("%x", md5.Sum(key))
				entry = entry + `{ "index" : { "_index" : "commits", "_id" : "` + id + `"} }`
			} else {
				entry = entry + `{ "index" : { "_index" : "commits" } }`
			}

			entry = entry + "\n" + string(data) + "\n"
		} else {
			log.Println(err)
		}
	}
	e.DoPost(path, entry)
	return
}
