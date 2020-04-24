# approuvez

[![GoDoc](https://godoc.org/github.com/mvisonneau/approuvez?status.svg)](https://godoc.org/github.com/mvisonneau/approuvez)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/approuvez)](https://goreportcard.com/report/github.com/mvisonneau/approuvez)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/approuvez.svg)](https://hub.docker.com/r/mvisonneau/approuvez/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/approuvez/status.svg)](https://cloud.drone.io/mvisonneau/approuvez)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/approuvez/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/approuvez?branch=master)

> ✅ Obtain live confirmation from people

## Usage

```bash
NAME:
   approuvez - command line helper to obtain live confirmation from relevant people

USAGE:
   approuvez [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint endpoint, -e endpoint  websocket endpoint to connect in order to listen for events [$APPROUVEZ_ENDPOINT]
   --token token, -t token           token to use in order to authenticate against the endpoint [$APPROUVEZ_TOKEN]
   --slack-token token               token to use in order to authenticate requests against slack [$APPROUVEZ_SLACK_TOKEN]
   --slack-message message           message to display on Slack [$APPROUVEZ_SLACK_MESSAGE]
   --slack-channel channel           slack channel to write the message onto [$APPROUVEZ_SLACK_CHANNEL]
   --triggerrer email or slack ID    email or slack ID of the person trigerring the message [$APPROUVEZ_TRIGGERRER]
   --reviewer email or slack ID      email or slack ID of a person that should review the message [$APPROUVEZ_REVIEWER]
   --required-approvals amount       amount of approvals required to consider it approved (default to all defined reviewers) (default: 0) [$APPROUVEZ_REQUIRED_APPROVALS]
   --log-level level                 log level (debug,info,warn,fatal,panic) (default: "info") [$APPROUVEZ_LOG_LEVEL]
   --log-format format               log format (json,text) (default: "text") [$APPROUVEZ_LOG_FORMAT]
   --help, -h                        show help
   --version, -v                     print the version
```

## Architecture

![approuvez_architecture](docs/images/approuvez_architecture.png)

## Develop / Test

```bash
~$ make build-local
~$ ./approuvez
```

## Build / Release

If you want to build and/or release your own version of `approuvez`, you need the following prerequisites :

- [git](https://git-scm.com/)
- [golang](https://golang.org/)
- [make](https://www.gnu.org/software/make/)
- [goreleaser](https://goreleaser.com/)

```bash
~$ git clone git@github.com:mvisonneau/approuvez.git && cd approuvez

# Build the binaries locally
~$ make build-local

# Build the binaries and release them (you will need a GITHUB_TOKEN and to reconfigure .goreleaser.yml)
~$ make release
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/approuvez/pulls).

## Fun fact

`approuvez` is a conjugation of the verb [approuver](https://www.larousse.fr/conjugaison/francais/approuver/518) in French 🇫🇷, equivalent to `approve` in English 🇬🇧
