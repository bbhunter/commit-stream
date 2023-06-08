package database

type Database interface {
	Insert(CommitEntry)
	Connect() error
}
