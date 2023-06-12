package database

import (
	"log"
	"reflect"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type Database interface {
	Insert(CommitEntry)
	Connect() error
}

type DatabaseHandler struct {
	Db Database
}

func NewDatabaseHandler(db Database) *DatabaseHandler {
	log.Println("Using database handler:" + reflect.TypeOf(db).String())
	databaseHandler := &DatabaseHandler{
		Db: db,
	}
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return databaseHandler

}
func (dh *DatabaseHandler) Callback(commits []commit.CommitEvent) {
	for _, c := range commits {
		email := c.AuthorEmail.User + "@" + c.AuthorEmail.Domain
		d := CommitEntry{
			AuthorName:  c.AuthorName,
			RepoName:    c.RepoName,
			UserName:    c.UserName,
			AuthorEmail: email,
			Message:     c.Message,
			SHA:         c.SHA,
		}
		dh.Db.Insert(d)
	}
}
