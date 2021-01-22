# ✅ approuvez

[![GoDoc](https://godoc.org/github.com/mvisonneau/approuvez?status.svg)](https://godoc.org/github.com/mvisonneau/approuvez)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/approuvez)](https://goreportcard.com/report/github.com/mvisonneau/approuvez)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/approuvez.svg)](https://hub.docker.com/r/mvisonneau/approuvez/)
[![Build Status](https://github.com/mvisonneau/approuvez/workflows/test/badge.svg?branch=main)](https://github.com/mvisonneau/approuvez/actions)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/approuvez/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/approuvez?branch=master)

## Usage

```bash
~$ approuvez --help
NAME:
   approuvez - Command line helper to obtain live confirmation from people in a blocking fashion

USAGE:
   approuvez [global options] command [command options] [arguments...]

COMMANDS:
   ask      send a message to someone and wait for a response
   serve    run the server thing
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level level    log level (debug,info,warn,fatal,panic) (default: "info") [$APPROUVEZ_LOG_LEVEL]
   --log-format format  log format (json,text) (default: "text") [$APPROUVEZ_LOG_FORMAT]
   --tls-disable        disable mutual tls for gRPC transmissions (use with care!) (default: false) [$APPROUVEZ_TLS_DISABLE]
   --tls-ca-cert path   TLS CA certificate path [$APPROUVEZ_TLS_CA_CERT]
   --tls-cert path      TLS certificate path [$APPROUVEZ_TLS_CERT]
   --tls-key path       TLS key path [$APPROUVEZ_TLS_KEY]
   --help, -h           show help (default: false)
```

### Server

```bash
~$ approuvez serve --help
NAME:
   approuvez serve - run the server thing

USAGE:
   approuvez serve [command options] [arguments...]

OPTIONS:
   --slack-token token     token to use in order to authenticate requests against slack [$APPROUVEZ_SLACK_TOKEN]
   --listen-address token  token to use in order to authenticate requests against slack (default: ":8443") [$APPROUVEZ_LISTEN_ADDRESS]
   --help, -h              show help (default: false)
```

### Client

```bash
~$ approuvez ask --help
NAME:
   approuvez ask - send a message to someone and wait for a response

USAGE:
   approuvez ask [command options] [arguments...]

OPTIONS:
   --endpoint endpoint, -e endpoint                server endpoint to connect upon (default: "127.0.0.1:8443") [$APPROUVEZ_SERVER_ENDPOINT]
   --user email or slack ID, -u email or slack ID  email or slack ID of a person that should review the message [$APPROUVEZ_USER]
   --message message, -m message                   message to display on Slack [$APPROUVEZ_MESSAGE]
   --link-name name                                name of a link button to append to the message [$APPROUVEZ_LINK_NAME]
   --link-url url                                  url of a link button to append to the message [$APPROUVEZ_LINK_URL]
   --help, -h                                      show help (default: false)
```

## Architecture

![approuvez_architecture](docs/images/approuvez_architecture.png)

## Develop / Test

```bash
~$ make build
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

## Terminology

`approuvez` is a conjugation of the verb [approuver](https://www.larousse.fr/conjugaison/francais/approuver/518) in French 🇫🇷, equivalent to `approve` in English 🇬🇧
