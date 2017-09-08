# motkid

The Moritake Industries EVE Corporation suite of applications.

[![GoDoc](https://godoc.org/github.com/motki/motkid?status.svg)](https://godoc.org/github.com/motki/motkid)

## Getting started

This repository contains all the motkid source code, including dependency source code.

#### Pre-requisites for building

* [A recent Go compiler](https://golang.org)
* [go-bindata](https://github.com/jteeuwen/go-bindata)
  ```bash
  go get -u github.com/jteeuwen/go-bindata/...
  ```
* postgres9 client (psql and pg_restore)
* cat, awk, sed, grep, curl, bunzip2, unzip, git

##### Requirements for running

* An available PostgreSQL database.
* An SMTP provider.
* A machine to run motkid on.


## Install using `make`

[Download a copy of this repository](https://github.com/motki/motkid/archive/master.zip) and verify that the Makefile works.

```bash
curl -L -o motkid.tar.gz https://github.com/motki/motkid/archive/master.tar.gz
tar xzf motkid.tar.gz
cd motkid
make debug
```

Assuming you haven't copied `config.toml.dist` to `config.toml`, you will be greeted with an error.

```
Makefile:61: *** config.toml does not exist. Copy config.toml.dist and edit appropriately, then try again..  Stop.
```

Once you've edited `config.toml`, you can actually build the program. The simplest way is `make`. After that, `make install`.


### Makefile reference

|  Target       | Description 
|-----------    |---------------------------------------------------
| build         | Build `motkid` and `motki`.
| install       | Installs database schemas and EVE static dump data.
| uninstall     | Drop created database schemas and EVE static dump data.
| clean         | Delete all build files.
| generate      | Runs `go generate`.
| matrix        | Build a matrix of arches and OSes, see below.
| download      | Download EVE static dump data.
| assets        | Installs EVE static dump data.
| db            | Install the database schemas.
| schema_evesde | Installs the EVE static dump schema.
| schema_app    | Installs the app schema.

#### Notes

Build for a specific OS and arch
```bash
make build GOOS=linux GOARCH=arm7
```

Cross-compile the binaries for many platforms at once
```bash
make matrix ARCHES="amd64 arm6 arm7 386" OSES="windows linux darwin"
```


## Manual Installation

Clone or [download the repository](https://github.com/motki/motkid/archive/master.zip).

This application does not rely on `$GOPATH`, but if you are planning on making changes, it may help to put it there.

Below is an example of one way to get the code.

> This assumes you have a simple `$GOPATH` with only one value (and no colons in it)

```bash
mkdir -p $GOPATH/src/github.com/motki
git clone git@github.com:motki/motkid $GOPATH/src/github.com/motki/motkid
cd $GOPATH/src/github.com/motki/motkid
```


### Install resources

Load the data in the `resources` folder.

1. Un-bzip the `evesde-*-postgres.dmp.bz2`.
2. Use `pg_restore` to load the EVE static dump.
   > Warnings abouts a missing "yaml" role can be ignored.
3. Extract the Icons and Types zips to `public/images` (creating `public/images/Icons` and `public/images/Types`)


## Configuration

Copy `config.toml.dist` to `config.toml` and edit appropriately.

#### Configuring the EVE API

To use the EVE API you need to set up an Application at the [EVE Developer Portal](https://developers.eveonline.com/applications).  You'll need to select appropriate roles (*all* of them is fine) and then set a correct Return URL for your setup.

> Note: the Return URL can include a port specification.

Once you have created your application on the developer portal, put the Client ID, Secret, and Return URL in the corresponding section in `config.toml`.


#### Configuring SSL

You need to generate a certificate and private key to properly set up SSL. During development, a self-signed certificate is recommended. For production deployments, the process is made simpler by using Let's Encrypt to automatically generate a valid certificate.


##### Generating a self-signed cert

1. Copy the source code from this stdlib utility: [generate_cert.go](https://golang.org/src/crypto/tls/generate_cert.go).
2. Put it inside its own package (something like `./cmd/gencert/` in the project directory).
3. Compile and run it: 
   `go run ./cmd/gencert/generate_cert.go --host localhost`
4. There should now be a `key.pem` and `cert.pem` file in the current working directory. Update `config.toml` with the path to these.
5. Start motkid

> Don't commit the `generate_cert` utility or the generated keys to the source code repository.

##### Generating a cert with letsencrypt

1. Configure the SSL section in `config.toml`
    1. Set `autocert=true` in config.toml.
    2. Set `certfile=""` and `keyfile=""` in config.toml
    3. Set the SSL `listen` parameter to a valid public hostname.
2. ...
3. Profit


### Running the app

For now, it's easiest to start `motkid` using `go run`.

```bash
cd $GOPATH/src/motki/motkid
go run ./cmd/motkid/main.go
```



