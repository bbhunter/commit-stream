package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// implements interface database
type Mysql struct {
	Dsn string
	db  *gorm.DB
}

func (d *Mysql) Connect() error {
	var err error
	d.db, err = gorm.Open(mysql.Open(d.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.db.AutoMigrate(&CommitEntry{})
	return nil
}

func (d *Mysql) Insert(c CommitEntry) {
	d.db.Create(&c)
}
