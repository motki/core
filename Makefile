# motkid Makefile

# This file acts as the heavy lifter for setting up an installation of motkid.
#
# Build both motki and motkid binaries:
#   make
#
# Download any necessary static data:
#   make download
#
# Install static assets and database schemas:
#   make install
#
# Remove static assets and database schemas:
#   make uninstall
#
# Build for a different OS and arch:
#   make build GOOS=linux GOARCH=arm7
#
# Cross-compile the binaries for many platforms at once:
#   make matrix RELEASE_ARCHES="amd64 arm6 arm7 386" RELEASE_OSES="windows linux darwin"
#
# Clean up build files:
#   make clean

# Defines where build files are stored.
PREFIX ?= build/

# Define all the necessary binary dependencies.
deps := go go-bindata psql pg_restore cat awk grep sed bunzip2 unzip git wget

# A template for defining a variable with the form:
#   NAME ?= /path/to/bin/name
deps_tpl = $(shell echo $(1) | tr a-z A-Z) ?= $(shell which $(1))

# Initialize a variable for each dependency listed in deps.
# The variables are upper-cased and overridable from the command-lane.
# For example, the variable containing the path to the go binary is called GO,
# psql is called PSQL, etc.
# Note that this variable's value is meaningless, it serves only as a name for this procedure.
deps_initialized := $(foreach dep,$(deps),$(eval $(call deps_tpl,$(dep))))

# Error messages.
err_go_missing := unable to locate "go" binary. See https://golang.org/doc/ for more information.
err_binary_missing = $(if $(filter $1,go),$(err_go_missing),unable to locate "$1" binary. \
Ensure that it is installed and on your PATH. Specify a custom path to the binary with \
"make $(shell echo $1 | tr a-z A-Z)=/path/to/$1 $@")
err_config_missing := config.toml does not exist. Copy config.toml.dist and edit appropriately, then try again.

# This procedure throws a fatal error if the path to the given binary is empty or
# does not exist.
ensure_dep = $(if $(realpath $(value $(shell echo $(1) | tr a-z A-Z))),,$(error $(call err_binary_missing,$(1))),exit 0;)

# Make sure we have all our dependencies. If a dependency is missing, make will
# exit with an appropriate error message.
# Note that this variable's value is meaningless, it serves only as a name for this procedure.
deps_ensured := $(foreach dep,$(deps),$(call ensure_dep,$(dep)))
# Make sure we have a configuration file to work with. If the file is missing, make will
# exit with an appropriate error message.
# Note that this variable's value is meaningless, it serves only as a name for this procedure.
config_ensured := $(if $(filter $(shell test -f config.toml && echo 1),1),,$(error $(err_config_missing)))


# Matches a URI: postgres://hostname:port/dbname?params
# Returns a sentence in the form: host:port dbname
db_conn_params := $(shell $(GREP) 'connection_string' ./config.toml \
	| $(SED) -E -e 's/^connection_string[\s]+?=[\s]+?\"postgres:\/?\/?(.*)\/(\w+)\??.*?\"$$/\1 \2/')
# Extract the host part of a host:port pair
host = $(word 1,$(subst :, ,$1))
# Extract the port part of a host:port pair.
# The second parameter functions as a default if a port is not specified.
port = $(or $(word 2,$(subst :, ,$1)),$(value 2))

# These fields define the PostgreSQL connection parameters.
# These are read from config.toml but can be customized when invoking make.
#   make db DB_HOST=example.com
DB_HOST ?= $(call host,$(word 1,$(db_conn_params)))
DB_PORT ?= $(call port,$(word 1,$(db_conn_params)),5432)
DB_NAME ?= $(word 2,$(db_conn_params))

# Arguments for invoking psql or pg_restore.
pg_args := -d $(DB_NAME) -h $(DB_HOST) -p $(DB_PORT)

# 1 if the schema $1 exists, otherwise 0.
schema_exists = $(PSQL) $(pg_args) -c "SELECT COUNT(1) FROM information_schema.schemata WHERE schema_name = '$1';" \
	| $(AWK) '{ if ($$1 == "1") { print "1" } if ($$1 == "0") { print "0" } }'
drop_schema = $(PSQL) $(pg_args) -c "DROP SCHEMA $1 CASCADE;" || exit 0


# By default, the system defined GOOS and GOARCH are used.
# These are overridable from the command line. For example:
#   make build GOOS=linux GOARCH=arm7
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

# When using the "matrix" target, these specify which OSes and archs to target.
# These are both overridable from the command line. For example:
#   make matrix RELEASE_ARCHES=amd64 arm6 arm7 386 RELEASE_OSES=linux
RELEASE_ARCHES ?= amd64
RELEASE_OSES ?= linux darwin

# Components to build up a valid "go build" command.
build_base := $(GO) build -ldflags "-s -w -X main.Version=$(shell $(GIT) describe --always)"
build_name = $(PREFIX)$1$(if $(filter $(GOOS),windows),.exe,)
build_src = ./cmd/$(word 1,$(subst _, ,$(subst ., ,$(subst $(PREFIX),,$1))))/*.go
build_cmd = GOOS=$(GOOS) GOARCH=$(GOARCH) $(build_base) -o $1 $(call build_src,$1)
release_name = $(call build_name,$1_$(GOOS)_$(GOARCH))
release_cmd = $(subst  build , build -tags release ,$(call build_cmd,$1))

# These define the programs configured. Adding more targets is
# automatic as long as the source code for the target exists in
# ./cmd/<target>/*.go.
binaries := motki motkid
binary_targets := $(foreach bin,$(binaries),$(call build_name,$(bin)))
release_targets := $(foreach bin,$(binaries),$(call release_name,$(bin)))

# These define the schema names.
# Note that targets for schemas are manually defined.
schemas := evesde app
schema_targets := $(foreach sch,$(schemas),schema_$(sch))

# These define where working EVE Static Dump data can be downloaded.
static_base_url := https://github.com/motki/motkid/raw/master/
download_targets := resources/evesde-postgres.dmp.bz2 resources/Icons.zip resources/Types.zip

# Static asset targets. These must be zip files and follow a specific
# convention. It may not be suitable for all assets.
assets := Types Icons
asset_images_dir := public/images/
asset_targets := $(foreach a,$(assets),$(asset_images_dir)$(a))


# Print configuration information: paths, build options, and config params.
extra_params := DB_HOST DB_PORT DB_NAME
define print_conf
	@$(foreach dep,$(deps),echo "$(shell echo $(dep) | tr a-z A-Z)=$(value $(shell echo $(dep) | tr a-z A-Z))";)
	@echo "GOOS=$(GOOS)"
	@echo "GOARCH=$(GOARCH)"
	@$(foreach val,$(extra_params),echo "$(val)=$($(val))";)
endef

# All of the files this generates.
files := $(PREFIX)motki_*_* $(PREFIX)motkid_*_* $(binary_targets)

.PHONY: all
.PHONY: build matrix release
.PHONY: install uninstall
.PHONY: download assets
.PHONY: db $(schema_targets)
.PHONY: clean clean_files
.PHONY: drop_schemas delete_assets
.PHONY: debug


# Build all binaries.
build: $(binary_targets)

# Install assets and initialize database schemas.
install: db assets


# This defines a target that matches any of the values listed in binary_targets.
$(binary_targets):
	$(call build_cmd,$@)

# Make release builds for the specified OS and arch.
release: $(release_targets)

$(release_targets):
	$(GO) generate
	$(call release_cmd,$@)

# This target will build a binary for every combination of
# RELEASE_ARCHES and RELEASE_OSES specified.
matrix:
	@for arch in $(RELEASE_ARCHES); do               \
		for os in $(RELEASE_OSES); do                \
			echo "Building $$os $$arch...";          \
			$(MAKE) release GOOS=$$os GOARCH=$$arch; \
		done;                                        \
	done;                                            \
	echo "Done."


# Installs the database schemas and data.
db: $(schema_targets)

# Installs the EVE static dump data if it does not already exist.
# Note that the pg_restore command will always exit with a zero.
# This is because the dump currently causes warnings to be emitted
# and pg_restore exits with an error exit code.
schema_evesde:
ifeq ($(shell echo $(call schema_exists,evesde)),0)
	($(BUNZIP2) -ck ./resources/evesde-postgres.dmp.bz2 | $(PG_RESTORE) $(pg_args) --clean; exit 0)
endif

# Installs the app schema.
schema_app:
ifeq ($(shell $(call schema_exists,app)),0)
	$(CAT) $(wildcard ./resources/ddl/*.sql) | $(PSQL) $(pg_args)
endif

# Downloads all necessary EVE Static Dump files.
download: $(download_targets)

# This defines a target that matches any static files that need to be downloaded.
$(download_targets):
	$(WGET) -P resources $(static_base_url)$@

# Installs all asset targets.
assets: download $(asset_targets)

# This defines a target for each asset defined in asset_targets.
# Note that it currently only supports the current structure with little
# flexibility.
$(asset_targets):
	$(UNZIP) ./resources/$(lastword $(subst /, ,$@)).zip -d $(asset_images_dir)

# Deletes build files and database schemas.
clean: clean_files

# Deletes all build files.
clean_files:
	@for f in $(files); do (rm -r "$$f" 2> /dev/null && echo "Deleted $$f"; exit 0); done
	@echo "Cleaned files."

# Deletes all schemas and removes installed assets.
uninstall: drop_schemas delete_assets clean_files

# Deletes all installed assets.
delete_assets:
	@for f in $(asset_targets); do (rm -r "$$f" 2> /dev/null && echo "Deleted $$f"; exit 0); done
	@echo "Deleted static assets."

# Deletes all database schemas.
drop_schemas:
	@$(foreach sch,$(schemas),$(if $(filter $(shell echo $(call schema_exists,$(sch))),1),$(call drop_schema,$(sch)) && echo "Dropped schema $(DB_NAME).$(sch)";,))
	@echo "Cleaned schemas."


# Prints configuration information.
debug:
	$(print_conf)
