package database

import (
	"gorm.io/gorm"
)

type CommitEntry struct {
	gorm.Model
	RepoName    string
	UserName    string
	AuthorName  string
	AuthorEmail string
	Message     string
	SHA         string
}
