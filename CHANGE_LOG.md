# Change Log

## [v1.11](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.11) (2021-08-18)

- sharding links across multiple servers (optionally):
  - supporting individual data storages for each server:
    - [Redis](https://redis.io/) database;
    - [MongoDB](https://www.mongodb.com/) database;
  - supporting specifying of the source server in a link:
    - returning the server ID on link generating;
    - ignoring the server ID:
      - on link getting;
      - on redirecting to a link.

## [v1.10](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.10) (2021-05-26)

- refactoring:
  - simplify the `handlers` package:
    - move to the `gateways` package;
    - merge with the `router` package;
    - simplify the `handlers.NewRouter()` function;
  - update the [github.com/thewizardplusplus/go-http-utils](https://github.com/thewizardplusplus/go-http-utils) package:
    - use the `httputils.ReadJSON()` function;
    - use the `httputils.WriteJSON()` function;
    - use the `httputils.ParsePathParameter()` function;
- distributing:
  - use the [wait-for-it.sh](https://github.com/vishnubob/wait-for-it) script to wait for dependencies to become available;
  - complete the [Docker Compose](https://docs.docker.com/compose/) configuration:
    - add the [Redis Commander](https://github.com/joeferner/redis-commander) service;
    - add the [mongo-express](https://github.com/mongo-express/mongo-express) service;
    - add the [e3w](https://github.com/soyking/e3w) service;
- add the API description in the format of a [Postman](https://www.postman.com/) collection.

## [v1.9](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.9) (2020-09-16)

- storing links in the [MongoDB](https://www.mongodb.com/) database:
  - create indexes for all link fields on connecting to the database;
  - don't insert a link but update in the upsert mode to avoid duplicates;
- refactoring:
  - of the `storage` package;
  - of integration tests.

## [v1.8](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.8) (2020-06-06)

- adding transformations of values received from a distributed counter:
  - using linear transformations:
    - multiplication;
    - shift;
- removing:
  - generating link codes:
    - removing using a [version 4 UUID](<https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)>);
    - using sequential counters:
      - removing formatting a counter as an integer number in the 10 base;
- refactoring:
  - extracting the `counters` package:
    - moving the `generators.chunkedCounter` structure to the `counters.ChunkedCounter` structure;
    - extracting the `counters.CounterGroup` structure from the `generators.DistributedGenerator` structure;
  - using external packages:
    - replacing `Printer` interfaces to the [special](https://github.com/go-log/log) package;
    - extracting HTTP utility functions to the [single](https://github.com/thewizardplusplus/go-http-utils) package;
- testing:
  - bulky testing:
    - adding a bulky test of the `generators.DistributedGenerator` structure;
    - adding an integration bulky test of link creating;
  - Travis CI configuration:
    - using data race detection;
    - using the [dep](https://golang.github.io/dep/) tool;
- improving repository decor:
  - adding the change log.

## [v1.7](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.7) (2020-03-10)

- renaming the project from "go-link-shortener" to "go-link-shortener-backend";
- improving:
  - project structure;
  - repository decor.

## [v1.6](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.6) (2020-03-05)

- removing the custom "Not Found" handler:
  - from the RESTful API;
  - from the additional routing;
- serving static files:
  - error handling:
    - tightening error handling when serving static files;
  - logging:
    - logging error handling when serving static files.

## [v1.5](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.5) (2020-02-29)

- generating link codes:
  - formatting a counter as an integer number in the 62 base.

## [v1.4](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.4) (2020-02-26)

- storing counters chunks in the [etcd](https://etcd.io/) database:
  - using a record version as a counter chunk instead of its revision.

## [v1.3](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.3) (2020-02-21)

- additional routing:
  - redirecting to the link URL by its code;
  - serving static files;
- passing a request to presenters:
  - of links;
  - of errors.

## [v1.2](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.2) (2019-12-20)

- error handling:
  - weakening caching error handling;
  - tightening response presenting error handling:
    - with links;
    - with errors;
- logging:
  - logging requests;
  - logging caching error handling;
  - logging response presenting error handling:
    - with links;
    - with errors;
- panics:
  - recovering on panics;
  - logging of panics;
- configuring time to live for link caches.

## [v1.1](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.1) (2019-11-08)

- generating link codes:
  - using sequential counters:
    - formatting:
      - formatting a counter as an integer number in the 10 base;
    - storing:
      - storing in a database only counters chunks;
      - storing counters themselves in memory;
    - sharding:
      - sharding counters chunks;
      - selecting a shard of a counter chunk at random;
    - using the [etcd](https://etcd.io/) database:
      - using a record revision as a counter chunk.

## [v1.0](https://github.com/thewizardplusplus/go-link-shortener-backend/tree/v1.0) (2019-10-18)
