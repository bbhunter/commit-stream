package main

import (
	"flag"

	commitstream "github.com/x1sec/commit-stream/pkg"
)

type FlagOptions struct {
	AuthToken  string
	Filter     commitstream.Filter
	ConfigFile string
	Debug      bool
}

func PopulateOptions(f *FlagOptions) {
	flag.StringVar(&f.Filter.Email, "email-domain", "", "")
	flag.StringVar(&f.Filter.Email, "e", "", "")

	flag.StringVar(&f.Filter.Name, "email-name", "", "")
	flag.StringVar(&f.Filter.Name, "n", "", "")

	flag.StringVar(&f.AuthToken, "token", "", "")
	flag.StringVar(&f.AuthToken, "t", "", "")

	flag.BoolVar(&f.Filter.IgnorePrivateEmails, "ignore-priv", false, "")
	flag.BoolVar(&f.Filter.IgnorePrivateEmails, "i", false, "")

	flag.BoolVar(&f.Filter.SearchAllCommits, "a", false, "")
	flag.BoolVar(&f.Filter.SearchAllCommits, "all-commits", false, "")
	flag.BoolVar(&f.Filter.IncludeMessages, "m", false, "")
	flag.BoolVar(&f.Filter.IncludeMessages, "messages", false, "")
	flag.StringVar(&f.ConfigFile, "config", "", "")
	flag.StringVar(&f.ConfigFile, "c", "", "")
	flag.BoolVar(&f.Debug, "debug", false, "")
	flag.BoolVar(&f.Debug, "d", false, "")
	flag.StringVar(&f.Filter.DomainsFile, "dom-file", "", "")
	flag.StringVar(&f.Filter.DomainsFile, "df", "", "")

	flag.BoolVar(&f.Filter.PublicEvent, "public-event", true, "")
	flag.BoolVar(&f.Filter.PublicEvent, "p", true, "")

	flag.Parse()
}
