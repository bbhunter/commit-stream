package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	SqLiteDB string
	db       *gorm.DB
}

func (d *Sqlite) Connect() error {
	var err error
	d.db, err = gorm.Open(sqlite.Open(d.SqLiteDB), &gorm.Config{})
	if err != nil {
		return err
	}
	d.db.AutoMigrate(&CommitEntry{})
	return nil
}

func (d *Sqlite) Insert(c CommitEntry) {
	d.db.Create(&c)
}

func (d *Sqlite) BatchInsert(c CommitEntry) {
	d.db.Create(c)
}
