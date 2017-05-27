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


### Install resources

Load the data in the `resources` folder.

* Use `pg_restore` to load the EVE static dump.
  > Warnings abouts a missing "yaml" role can be ignored.
* Extract the Icons, Renders, and Types zips to `public/images` (creating `public/images/Icons`, `public/images/Renders`, and `public/images/Types`)


### Configuration

Copy `config.toml.dist` to `config.toml` and edit appropriately.

#### Configuring the EVE API

To use the EVE API you need to set up an Application at the [EVE Developer Portal](https://developers.eveonline.com/applications).  You'll need to select appropriate roles (*all* of them is fine) and then set a correct Return URL for your setup.

> Note: the Return URL can include a port specification.

Once you have created your application on the developer portal, put the Client ID, Secret, and Return URL in the corresponding section in `config.toml`.

If you're using SSL, you currently must set `require=true` in the SSL configuration in `config.toml` and use the SSL listen address as the return URL for the EVE API Application.

#### Configuring SSL

You need to generate a certificate and private key to properly set up SSL. You can generate a self-signed cert or use something like Let's Encrypt.

If you're using SSL, you currently must set `require=true` in the SSL configuration in `config.toml` and use the SSL listen address as the return URL for the EVE API Application.

##### Generating a self-signed cert

1. Copy the source code from this stdlib utility: [generate_cert.go](https://golang.org/src/crypto/tls/generate_cert.go).
2. Put it inside its own package (something like `./cmd/gencert/` in the project directory).
3. Compile and run it: 
   `go run ./cmd/gencert/generate_cert.go --host localhost`
4. There should now be a `key.pem` and `cert.pem` file in the current working directory. Update `config.toml` with the path to these.
5. Start motkid

> Don't commit the `generate_cert` utility or the generated keys to the source code repository.

##### Generating a cert with letsencrypt

TBD


### Running the app

For now, it's easiest to start `motkid` using `go run`.

```bash
cd $GOPATH/src/tyler-sommer/motki
go run ./cmd/motkid/main.go
```


#### Building and deploying

TBD