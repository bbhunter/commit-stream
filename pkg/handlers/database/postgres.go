package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// implements interface database
type Postgres struct {
	Dsn string
	db  *gorm.DB
}

func (d *Postgres) Connect() error {
	var err error
	d.db, err = gorm.Open(postgres.Open(d.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.db.AutoMigrate(&CommitEntry{})
	return nil
}

func (d *Postgres) Insert(c CommitEntry) {
	d.db.Create(&c)
}
