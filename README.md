# MOTKI Core Libraries

The MOTKI EVE Corporation application libraries.

View [current documentation](https://godoc.org/github.com/motki/core) on GoDoc.

[![Build Status](https://travis-ci.org/motki/core.svg?branch=master)](https://travis-ci.org/motki/core) [![GoDoc](https://godoc.org/github.com/motki/core?status.svg)](https://godoc.org/github.com/motki/core)

Subpackage Overview
-------------------

| Name               | Description
| ------------------ | -----------
| [app][1]           | Integration package that wires up all the dependencies in a MOTKI application. Use this package to bootstrap your own MOTKI client application.
| [db][2]            | PostgreSQL database integration. Light wrapper around [jackc/pgx](https://github.com/jackc/pgx).
| [eveapi][3]        | EVE API integration. Handles EVE SSO and fetching data from both ESI and XML APIs using [antihax/goesi](https://github.com/antihax/goesi).
| [evedb][4]         | EVE Static Data Export interface. Queries the SDE for static type/universe information. MOTKI uses [Fuzzwork's Postgres dump](https://www.fuzzwork.co.uk/dump/).
| [evemarketer][5]   | Provides region- and system-specific market statistics using [evemarketer.com](https://evemarketer.com).
| [log][6]           | Wrapper around [sirupsen/logrus](https://github.com/sirupsen/logrus) providing a configuration API and a defacto `Logger` type.
| [model][7]         | Encapsulates persistence of data to the database. General pattern is to fetch from DB, then from API if stale. The database schema for this package is defined in the [resources/ddl/ directory](https://github.com/motki/core/tree/master/resources/ddl).
| [proto][8]         | Defines the protocol buffer (and [gRPC](https://grpc.io)) interface for MOTKI at large.
| [proto/client][9]  | A golang gRPC client for interacting with a remote MOTKI application server.
| [proto/server][10] | A golang gRPC server for handling MOTKI client requests.

[1]: https://godoc.org/github.com/motki/core/app
[2]: https://godoc.org/github.com/motki/core/db
[3]: https://godoc.org/github.com/motki/core/eveapi
[4]: https://godoc.org/github.com/motki/core/evedb
[5]: https://godoc.org/github.com/motki/core/evemarketer
[6]: https://godoc.org/github.com/motki/core/log
[7]: https://godoc.org/github.com/motki/core/model
[8]: https://godoc.org/github.com/motki/core/proto
[9]: https://godoc.org/github.com/motki/core/proto/client
[10]: https://godoc.org/github.com/motki/core/proto/server
