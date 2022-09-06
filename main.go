/*
commit-stream
Author: https://twitter.com/robhax

See LICENSE
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	commitstream "github.com/robhax/commit-stream/pkg"
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
		h += "  -e, --email           Match email addresses field (specify multiple with comma). "
		h += "						  Omit to match all.\n"
		h += "  -n, --name            Match author name field (specify multiple with comma).\n"
		h += "                        Omit to match all.\n"
		h += "  -t, --token           Github token (if not specified, will use environment\n"
		h += "                        variable 'CSTREAM_TOKEN')\n"
		h += "  -a  --all-commits     Search through previous commit history (default: false)\n"
		h += "  -i  --ignore-priv     Ignore noreply.github.com private email addresses (default: false)\n"
		h += "  -m  --messages        Fetch commit messages (default: false)\n"
		h += "  -c  --config [ path ] Use configuration file (default: config.yaml)\n"
		h += "  -d  --debug           Enable debug messages (to stdout)"
		h += "\n\n"
		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {

	var (
		authToken  string
		filter     commitstream.Filter
		configFile string
		handler    commitstream.Handler
		debug      bool
	)

	flag.StringVar(&filter.Email, "email", "", "")
	flag.StringVar(&filter.Email, "e", "", "")

	flag.StringVar(&filter.Name, "name", "", "")
	flag.StringVar(&filter.Name, "n", "", "")

	flag.StringVar(&authToken, "token", "", "")
	flag.StringVar(&authToken, "t", "", "")

	flag.BoolVar(&filter.IgnorePrivateEmails, "ignore-priv", false, "")
	flag.BoolVar(&filter.IgnorePrivateEmails, "i", false, "")

	flag.BoolVar(&filter.SearchAllCommits, "a", false, "")
	flag.BoolVar(&filter.SearchAllCommits, "all-commits", false, "")
	flag.BoolVar(&filter.IncludeMessages, "m", false, "")
	flag.BoolVar(&filter.IncludeMessages, "messages", false, "")
	flag.StringVar(&configFile, "config", "", "")
	flag.StringVar(&configFile, "c", "", "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&debug, "d", false, "")

	flag.Parse()

	config := commitstream.Config{FilePath: configFile}
	if err := config.Load(); err != nil {
		log.Printf(err.Error())
	}
	fmt.Println(config.Settings.Destination)
	if filter.Email == "" && filter.Name == "" {
		filter.Enabled = false
	} else {
		filter.Enabled = true
	}

	if authToken == "" {
		if config.Settings.Github.Token == "" {
			authToken = os.Getenv("CSTREAM_TOKEN")
		} else {
			authToken = config.Settings.Github.Token
		}

		if authToken == "" {
			log.Fatal("No Github token specified. Use '-t', environment variable CSTREAM_TOKEN or specifying in config.yaml\n")
		}
	}

	if config.Settings.Destination == "elastic" {
		settings := config.Settings.Elasticsearch
		log.Printf("Using ElasticSearch database: %s\n", settings.Uri)
		handler = commitstream.ElasticHandler{
			RemoteURI: settings.Uri,
			Username:  settings.Username,
			Password:  settings.Password,
		}

	} else {
		log.Println("Outputting to stdout")
		handler = commitstream.CsvHander{}
	}

	githubOptions := commitstream.GithubOptions{
		AuthToken: authToken,
	}

	cs := commitstream.CommitStream{
		GithubOptions: &githubOptions,
		Filter:        &filter,
		Debug:         debug,
	}

	//handler = commitstream.NoHandler{}
	cs.Start(handler)
}
