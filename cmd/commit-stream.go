package main

import (
	"fmt"
	"flag"
	"os"
	"encoding/csv"
	"github.com/x1sec/commit-stream/pkg"
)


func init() {
	flag.Usage = func() {
		h := "Extract Github commit author details in realtime\n\n"

		h += "Usage:\n"
		h += "  commit-stream [OPTIONS]\n\n"

		h += "Options:\n"
		h += "  -e, --email       Match email addresses field (specify multiple with comma). Omit to match all.\n"
		h += "  -n, --name        Match author name field (specify multiple with comma). Omit to match all.\n"
		h += "  -t, --token       Github token (mandatory)\n"
		h += "  -o, --output      Fields to output. Supply one or more options:\n"
		h += "  -a  --all-commits Match against all commit events for each repository (default: false)\n"
		h += "                      d - date\n"
		h += "                      e - email address\n"
		h += "                      n - Author name\n"
		h += "                      r - Repository URL\n"
		h += "                    Example: '-o er' will output email and repository URL\n"
		h += "                    Default (enr)\n"


		fmt.Fprintf(os.Stderr, h)

	}
}

func main() {

	var (
		authToken   string
		outputFmt	string
		rate        int
		fo commitstream.FilterOptions
		fmtFlags		commitstream.FlagOption
		searchAllCommits  bool

	)
 

	flag.StringVar(&fo.Email, "email", "", "")
	flag.StringVar(&fo.Email, "e", "", "")

	flag.StringVar(&fo.Name, "name", "", "")
	flag.StringVar(&fo.Name, "n", "", "")

	flag.StringVar(&authToken, "token", "", "")
	flag.StringVar(&authToken, "t", "", "")	
	flag.IntVar(&rate, "r", 0,"")
	flag.IntVar(&rate, "rate", 0, "")

	flag.StringVar(&outputFmt, "output", "","")
	flag.StringVar(&outputFmt, "o", "","")

	flag.BoolVar(&searchAllCommits, "a", false, "")
	flag.BoolVar(&searchAllCommits, "all-commits", false, "")

	flag.Parse()

	commitstream.ParseOutputOption(&fmtFlags, outputFmt)


	if fo.Email == "" && fo.Name == "" {
		fo.Enabled = false
	} else {
		fo.Enabled = true
	}


	if authToken == "" {
		fmt.Fprintf(os.Stderr, "No auth token specified\n")
		os.Exit(1)
	}


	streamOpt := commitstream.StreamOptions{AuthToken : authToken, SearchAllCommits: searchAllCommits}
	commitstream.DoIngest(streamOpt, fo, fmtFlags, handleResult)

}

func handleResult(s []string) {
	w := csv.NewWriter(os.Stdout)
	w.Write(s)
	w.Flush()
}

