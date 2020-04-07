# Commit Stream

`commit-stream` drinks commit logs from the Github event firehose exposing the author details (name and email address) associated with Github repositories in real time. 

OSINT / Recon uses for  Redteamers / Bug bounty hunters: 

* Uncover repositories which employees of a target company is commiting code (filter by email domain)
* Identify repositories belonging to an individual (filter by author name)
* Chain with other tools such as trufflehog to extract secrets in uncovered repositories.

[![asciicast](https://asciinema.org/a/317469.svg)](https://asciinema.org/a/317469)

## Installation
### Binaries
Install the latest binaries from the releases or build yourself:

### Go get
If your go environment is setup:
```
go get -u github.com/x1sec/commit-stream
```
### Compiling
```
git clone https://github.com/x1sec/commit-stream
cd commit-stream
go get && go build
```

# Usage

```
Usage:
  commit-stream [OPTIONS]

Options:
  -e, --email       Match email addresses field (specify multiple with comma). Omit to match all.
  -n, --name        Match author name field (specify multiple with comma). Omit to match all.
  -t, --token       Github token (if not specified, will use environment variable 'CSTREAM_TOKEN')
  -a  --all-commits Search through previous commit history (default: false)
```

`commit-stream` requires a Github personal access token to be used. You can generate a token navigating in Github [Settings / Developer Settings /  Personal Access Tokens] then selecting 'Generate new token'. Nothing here needs to be selected, just enter the name of the token and click generate.

Once the token has been created, the recommended method is to set it via an environment variable `CSTREAM_TOKEN`:
```
export CSTREAM_TOKEN xxxxxxxxxx
```
Alternatively, the `--token` switch maybe used when invoking the program, e.g:
```
./commit-stream --token xxxxxxxxxx
```

When running `commit-stream` with no options, it will immediately dump author details and the associated repositories in CSV format to the terminal. Filtering options are available. 

To filter by email domain:
```
./commit-stream --email '@company.com'
```

To filter by author name:
```
./commit-stream --name 'John Smith'
```

Multiple keywords can be specified with a `,` character. e.g.
```
./commit-stream --email '@telsa.com,@ford.com`
```

It is possible to search upto 20 previous commits for the filter keywords by specifying `--all-commits`. This may increase the likelihood of a positive matches.

### Note
As only one token is used this software does not breach any terms of use with Github. That said, use at your own risk. The author does not hold any responsibility for it's usage.


