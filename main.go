/*
commit-stream
Author: https://twitter.com/robhax

See LICENSE
*/

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
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
		h += "  -e, --email        Match email addresses field (specify multiple with comma). Omit to match all.\n"
		h += "  -n, --name         Match author name field (specify multiple with comma). Omit to match all.\n"
		h += "  -t, --token        Github token (if not specified, will use environment variable 'CSTREAM_TOKEN')\n"
		h += "  -a  --all-commits  Search through previous commit history (default: false)\n"
		h += "  -i  --ignore-priv  Ignore noreply.github.com private email addresses (default: false)\n"
		h += "  -m  --messages     Fetch commit messages (default: false)\n"
		h += "\n\n"
		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {

	var (
		authToken string
		filter    commitstream.Filter
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

	flag.Parse()

	if filter.Email == "" && filter.Name == "" {
		filter.Enabled = false
	} else {
		filter.Enabled = true
	}

	if authToken == "" {
		authToken = os.Getenv("CSTREAM_TOKEN")
		if authToken == "" {
			fmt.Fprintf(os.Stderr, "Please specify Github authentication token with '-t' or by setting the environment variable CSTREAM_TOKEN\n")
			os.Exit(1)
		}
	}

	githubOptions := commitstream.GithubOptions{
		AuthToken: authToken,
	}

	cs := commitstream.CommitStream{
		GithubOptions: &githubOptions,
		Filter:        &filter,
	}

	cs.Start(handleResult)
}

func handleResult(c commitstream.Commit) {
	cOut := []string{c.Name, c.Email, "https://github.com/" + c.Repo, c.Message}
	w := csv.NewWriter(os.Stdout)
	w.Write(cOut)
	w.Flush()
}
