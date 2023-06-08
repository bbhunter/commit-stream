package filter

import (
	"strings"

	"github.com/x1sec/commit-stream/pkg/commit"
)

type Filter struct {
	Email                 string
	Name                  string
	Enabled               bool
	IgnorePrivateEmails   bool
	IncludeMessages       bool
	SearchAllCommitEvents bool
	DomainsFile           string
	DomainsList           map[string]bool
}

func Filtered(c commit.CommitEvent, filter Filter) bool {

	if filter.IgnorePrivateEmails == true {
		if strings.Contains(c.AuthorEmail.Domain, "users.noreply.github.com") {
			return false
		}
	}

	if filter.Enabled == false {
		return true
	}

	result := false

	if len(filter.DomainsList) != 0 {
		if ok := filter.DomainsList[c.AuthorEmail.Domain]; ok {
			return true
		}
	}
	if filter.Email != "" {
		for _, e := range strings.Split(filter.Email, ",") {
			if strings.Contains(c.AuthorEmail.Domain, strings.TrimSpace(e)) {
				result = true
			}
		}
	}

	if filter.Name != "" {
		for _, n := range strings.Split(filter.Name, ",") {
			if strings.Contains(c.AuthorName, strings.TrimSpace(n)) {
				result = true
			}
		}
	}

	return result
}
