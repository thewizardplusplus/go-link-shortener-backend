# go-link-shortener-backend

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-link-shortener-backend?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-link-shortener-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-link-shortener-backend)](https://goreportcard.com/report/github.com/thewizardplusplus/go-link-shortener-backend)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-link-shortener-backend.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-link-shortener-backend)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-link-shortener-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-link-shortener-backend)

Back-end of the service for shorting links.

## Features

- RESTful API:
  - link model:
    - creating by an URL;
    - getting by a code;
  - representing in a JSON:
    - payloads:
      - of requests;
      - of responses;
    - errors;
- generating link codes:
  - using sequential counters:
    - formatting:
      - formatting a counter as an integer number in the 62 base;
    - storing:
      - storing in a database only counters chunks;
      - storing counters themselves in memory;
    - sharding:
      - sharding counters chunks;
      - selecting a shard of a counter chunk at random;
- server:
  - additional routing:
    - redirecting to the link URL by its code;
    - serving static files;
  - storing settings in environment variables;
  - supporting graceful shutdown;
  - logging:
    - logging requests;
    - logging errors;
  - panics:
    - recovering on panics;
    - logging of panics;
- databases:
  - storing links in the [MongoDB](https://www.mongodb.com/) database;
  - storing counters chunks in the [etcd](https://etcd.io/) database:
    - using a record version as a counter chunk;
  - caching links in the [Redis](https://redis.io/) database;
- distributing:
  - [Docker](https://www.docker.com/) image;
  - [Docker Compose](https://docs.docker.com/compose/) configuration.

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

API description:

- in the [Swagger](http://swagger.io/) format: [docs/swagger.yaml](docs/swagger.yaml);
- in the format of a [Postman](https://www.postman.com/) collection: [docs/postman_collection.json](docs/postman_collection.json).

## Bibliography

1. [URL Shortener System Design](https://medium.com/@narengowda/url-shortener-system-design-3db520939a1c).
2. [Generating Globally Unique Identifiers for Use with MongoDB](https://www.mongodb.com/blog/post/generating-globally-unique-identifiers-for-use-with-mongodb).

## License

The MIT License (MIT)

Copyright &copy; 2019-2020 thewizardplusplus
