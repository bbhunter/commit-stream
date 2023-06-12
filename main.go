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

	githubOptions := github.GithubOptions{
		AuthToken: authToken,
	}

	cs := stream.CommitEventStream{
		GithubOptions: &githubOptions,
		Filter:        &flags.Filter,
		Debug:         flags.Debug,
	}

	handler = newHandler(config)
	cs.Start(handler)
}
