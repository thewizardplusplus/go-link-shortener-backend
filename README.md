# go-link-shortener-backend

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-link-shortener-backend?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-link-shortener-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-link-shortener-backend)](https://goreportcard.com/report/github.com/thewizardplusplus/go-link-shortener-backend)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-link-shortener-backend.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-link-shortener-backend)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-link-shortener-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-link-shortener-backend)

Back-end of the service for shorting links.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-link-shortener-backend.git
$ cd go-link-shortener-backend
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

Build the project:

```
$ go install ./...
```

## Usage

```
$ go-link-shortener
```

Environment variables:

- `SERVER_ADDRESS` &mdash; server URI (default: `:8080`);
- `SERVER_STATIC_PATH` &mdash; path to the project's front-end (default: `./static`);
- `CACHE_ADDRESS` &mdash; [Redis](https://redis.io/) connection URI (default: `localhost:6379`);
- `CACHE_TTL_CODE` &mdash; time to live of links in [Redis](https://redis.io/), stored by their code (e.g. `72h3m0.5s`; default: `1h`);
- `CACHE_TTL_URL` &mdash; time to live of links in [Redis](https://redis.io/), stored by their URL (e.g. `72h3m0.5s`; default: `1h`);
- `STORAGE_ADDRESS` &mdash; [MongoDB](https://www.mongodb.com/) connection URI (default: `mongodb://localhost:27017`);
- `COUNTER_ADDRESS` &mdash; [etcd](https://etcd.io/) connection URI (default: `localhost:2379`);
- `COUNTER_COUNT` &mdash; count of distributed counters (default: `2`);
- `COUNTER_CHUNK` &mdash; step of a distributed counter (default: `1000`);
- `COUNTER_RANGE` &mdash; range of a distributed counter (default: `1000000000`).

## API Description

API description in the [Swagger](http://swagger.io/) format: [docs/swagger.yaml](docs/swagger.yaml).

## Bibliography

1. [URL shortener System design](https://medium.com/@narengowda/url-shortener-system-design-3db520939a1c).
2. [Generating Globally Unique Identifiers for Use with MongoDB](https://www.mongodb.com/blog/post/generating-globally-unique-identifiers-for-use-with-mongodb).

## License

The MIT License (MIT)

Copyright &copy; 2019-2020 thewizardplusplus
