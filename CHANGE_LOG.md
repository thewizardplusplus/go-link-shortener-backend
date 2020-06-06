# Change Log

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
