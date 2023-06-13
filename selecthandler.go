package main

import (
	"log"
	"os"

	"github.com/x1sec/commit-stream/pkg/conf"
	"github.com/x1sec/commit-stream/pkg/handlers"
	"github.com/x1sec/commit-stream/pkg/handlers/csv"
	"github.com/x1sec/commit-stream/pkg/handlers/database"
	"github.com/x1sec/commit-stream/pkg/handlers/elastic"
	"github.com/x1sec/commit-stream/pkg/handlers/script"
	"github.com/x1sec/commit-stream/pkg/handlers/slack"
	"github.com/x1sec/commit-stream/pkg/handlers/truffle"
)

func newHandler(config conf.Config) handlers.Handler {
	s := config.Settings
	switch s.Destination {
	case "":
		return csv.NewCsvHandler(os.Stdout)
	case "csv":
		return csv.NewCsvHandler(os.Stdout)
	case "elastic":
		return elastic.NewElasticHandler(s.Elasticsearch.Uri,
			s.Elasticsearch.Username,
			s.Elasticsearch.Password,
			s.Elasticsearch.NoDuplicates)
	case "script":
		return script.NewScriptHandler(s.Script.Path,
			s.Script.MaxWorkers,
			s.Script.LogFile)
	case "truffle-slack":
		slackConf := slack.NewSlackHandler(s.Slack.Token, s.Slack.ChannelID)
		ts := truffle.NewTruffleHandler(
			s.Truffle.Path,
			s.Truffle.MaxWorkers,
			s.Truffle.GitHubToken)
		ts.SlackConf = slackConf
		log.Println("Truffle matches to slack")
		return ts
	case "slack":
		return slack.NewSlackHandler(s.Slack.Token,
			s.Slack.ChannelID)
	case "database":
		db := selectDatabaseHandler(config)
		return database.NewDatabaseHandler(db)
	default:
		log.Fatal("Unknown handler: " + s.Destination)
	}
	return nil
}
func selectDatabaseHandler(config conf.Config) database.Database {
	var db database.Database
	engine := config.Settings.Database.Engine
	if engine == "sqlite" {
		path := config.Settings.Database.Path
		db = &database.Sqlite{
			SqLiteDB: path,
		}

	} else if engine == "postgres" {
		dsn := config.Settings.Database.Dsn
		db = &database.Postgres{
			Dsn: dsn,
		}
	} else if engine == "mysql" {
		dsn := config.Settings.Database.Dsn
		db = &database.Mysql{
			Dsn: dsn,
		}
	} else {
		log.Fatal("Unknown database engine: " + engine)
	}
	return db
}
