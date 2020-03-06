# ✅ approuvez - obtain live confirmation from people

[![GoDoc](https://godoc.org/github.com/mvisonneau/approuvez?status.svg)](https://godoc.org/github.com/mvisonneau/approuvez)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/approuvez)](https://goreportcard.com/report/github.com/mvisonneau/approuvez)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/approuvez.svg)](https://hub.docker.com/r/mvisonneau/approuvez/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/approuvez/status.svg)](https://cloud.drone.io/mvisonneau/approuvez)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/approuvez/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/approuvez?branch=master)

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
