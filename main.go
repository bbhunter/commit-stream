/*
commit-stream
Author: https://twitter.com/x1sec

See LICENSE
*/

package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/x1sec/commit-stream/pkg/conf"
	"github.com/x1sec/commit-stream/pkg/database"
	"github.com/x1sec/commit-stream/pkg/github"
	"github.com/x1sec/commit-stream/pkg/handlers"
	"github.com/x1sec/commit-stream/pkg/stream"
)

func printAscii() {
	h := `
 ██████╗ ██████╗ ███╗   ███╗███╗   ███╗██╗████████╗   ███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗
██╔════╝██╔═══██╗████╗ ████║████╗ ████║██║╚══██╔══╝   ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║
██║     ██║   ██║██╔████╔██║██╔████╔██║██║   ██║█████╗███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║
██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██║   ██║╚════╝╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║
╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║   ██║      ███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║
 ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝   ╚═╝      ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝ 
https://github.com/x1sec/commit-stream       

`
	fmt.Fprintf(os.Stderr, h)
}

func init() {
	flag.Usage = func() {
		printAscii()

		h := "Stream Github commit logs in real-time\n\n"
		h += "Usage:\n"
		h += "  commit-stream [OPTIONS]\n\n"

		h += "Options:\n"
		h += "  -t, --token            Github token (if not specified, will use environment\n"
		h += "                         variable 'CSTREAM_TOKEN' or from config.yaml)\n"
		h += "  -e, --email-domain     Match email addresses field (specify multiple with comma)\n"
		h += "                         Omit to match all.\n"
		h += "  -n, --email-name       Match author name field (specify multiple with comma).\n"
		h += "                         Omit to match all.\n"
		h += "  -df --dom-file <file>  Match email domains specificed in file\n"
		h += "  -a  --all-commits      Search through previous commit history (default: false)\n"
		h += "  -i  --ignore-priv      Ignore noreply.github.com private email addresses (default: false)\n"
		h += "  -m  --messages         Fetch commit messages (default: false)\n"
		h += "  -c  --config [path]    Use configuration file. Required for ElasticSearch (default: config.yaml)\n"
		h += "  -d  --debug            Enable debug messages to stderr (default:false)\n"
		h += "  -h  --help             This message\n"
		h += "\n\n"
		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {
	var handler handlers.Handler
	var flags conf.FlagOptions
	conf.PopulateOptions(&flags)

	var authToken string
	config := conf.Config{FilePath: flags.ConfigFile}
	if err := config.Load(); err != nil {
		log.Printf(err.Error())
	}

	var level log.Level
	if config.Settings.LogLevel == "debug" {
		log.Println("log level debug")
		level = log.DebugLevel
	} else {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	if flags.AuthToken == "" {
		if config.Settings.Github.Token == "" {
			authToken = os.Getenv("CSTREAM_TOKEN")
		} else {
			authToken = config.Settings.Github.Token
		}

	} else {
		authToken = flags.AuthToken
	}
	if authToken == "" {
		log.Fatal("No Github token specified. Use '-t', environment variable CSTREAM_TOKEN or specifying in config.yaml\n")
	}

	if config.Settings.Destination == "elastic" {
		settings := config.Settings.Elasticsearch

		log.Printf("Using ElasticSearch database: %s\n", settings.Uri)
		handler = handlers.ElasticHandler{
			RemoteURI:    settings.Uri,
			Username:     settings.Username,
			Password:     settings.Password,
			NoDuplicates: settings.NoDuplicates,
			UseZincAwsS3: settings.UseZincAwsS3,
		}
		h := handler.(handlers.ElasticHandler)
		h.Setup()

	} else if config.Settings.Destination == "script" {
		log.Println("Using script handler: " + config.Settings.Script.Path)

		handler = handlers.ScriptHandler{
			Path:       config.Settings.Script.Path,
			MaxWorkers: config.Settings.Script.MaxWorkers,
			LogFile:    config.Settings.Script.LogFile,
		}
	} else if config.Settings.Destination == "truffle-slack" {
		log.Println("Running trufflehog on commits and sending to slack")
		token := config.Settings.Slack.Token
		channelID := config.Settings.Slack.ChannelID
		slack := handlers.NewSlack(token, channelID)
		h := handlers.TruffleHandler{
			Slack:       slack,
			Path:        config.Settings.Truffle.Path,
			MaxWorkers:  config.Settings.Truffle.MaxWorkers,
			GithubToken: config.Settings.Truffle.GitHubToken,
		}
		h.StartWorkers()
		handler = &h
	} else if config.Settings.Destination == "slack" {
		log.Println("Using slack handler")
		if flags.Filter.Email == "" && flags.Filter.Name == "" {
			log.Fatal("No filter options specified. Refusing to use slack handler!")
		}
		token := config.Settings.Slack.Token
		channelID := config.Settings.Slack.ChannelID
		handler = handlers.NewSlack(token, channelID)

	} else if config.Settings.Destination == "database" {
		log.Println("Using database handler")
		var databaseHandler handlers.DatabaseHandler
		var db database.Database
		engine := config.Settings.Database.Engine
		log.Println("\t.. engine: " + engine)
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
			log.Fatal("Unknown database engine. Exiting.")
		}
		databaseHandler = handlers.DatabaseHandler{
			Db: db,
		}
		err := db.Connect()
		if err != nil {
			log.Fatal(err)
		}
		handler = &databaseHandler
	} else {
		handler = handlers.CsvHander{}
	}

	githubOptions := github.GithubOptions{
		AuthToken: authToken,
	}

	cs := stream.CommitEventStream{
		GithubOptions: &githubOptions,
		Filter:        &flags.Filter,
		Debug:         flags.Debug,
	}

	cs.Start(handler)
}
