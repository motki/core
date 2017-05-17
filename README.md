# MOTKI

Moritake Industries EVE Corporation website


## Installation

Clone the repository to the appropriate place in your `$GOPATH`.

> This assumes you have a simple `$GOPATH` with only one value (and no colons in it)

```bash
mkdir -p $GOPATH/src/github.com/tyler-sommer
git clone git@github.com:tyler-sommer/motki $GOPATH/src/tyler-sommer/motki
cd $GOPATH/src/tyler-sommer/motki
```

#### Install resources

Load the data in the resources folder.

* Use `pg_restore` to load the EVE static dump.
  > Warnings abouts a missing "yaml" role can be ignored.
* Extract the Icons, Renders, and Types zips to `public/images` (creating `public/images/Icons`, `public/images/Renders`, and `public/images/Types`)

#### Configuration

Copy `config.toml.dist` to `config.toml` and edit appropriately.

#### Running the app

For now, it's easiest to start `motkid` using `go run`.

```bash
cd $GOPATH/src/tyler-sommer/motki
go run ./cmd/motkid/main.go
```