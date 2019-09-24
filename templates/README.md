# kumgo-skelly
Kumparan Framework for Golang

## Project Structure
1. **pb** - This folder is used to store the *.proto files and the compilation result.
2. **client** - Client folder is used for other services clients.
3. **config** - This folder used for the configuration setting.
4. **console** - All the terminal commands goes here.
5. **deploy** - Deployment scripts.
6. **service** - Service implementation.
7. **repository** - Service model, query and logic.
8. **db** - DB connection, migration.

## Installing Go
### MacOS
1. You can follow the steps [here](https://golang.org/doc/install#osx)
2. Make sure you set your `$GOPATH` right. Check how to set your `$GOPATH` [here](https://github.com/golang/go/wiki/SettingGOPATH)

Some plugin will need to be installed in the `$GOBIN`. Default `$GOBIN` is `$GOPATH/bin`. It must be in your `$PATH` for the plugin to find it.
1. Add this line in your .zshrc or .bash_profile
```shell
$ export PATH=$PATH:$GOPATH/bin
```
2. Then reload the profile
```shell
$ source ~/.zshrc #if you are using zsh
$ source ~/.bash_profile #if you are using bash
```

## Installing Protobuf Compiler
### MacOS
1. Download the appropriate release [here](https://github.com/google/protobuf/releases)
2. Unzip the folder
3. Enter the folder and run `./autogen.sh && ./configure && make`
...If you run into this `error: autoreconf: failed to run aclocal: No such file or directory`, run `brew install libtool && brew install autoconf && brew install automake`. And run the command from step 3 again.
4. Then run these other commands. They should run without issues
```shell
$ make check
$ sudo make install
$ which protoc
$ protoc --version
```
5. Install the protoc plugin for go.
```shell
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

## Installing Dep (Go Dependency Manager)
You can follow the steps [here](https://github.com/golang/dep)

## Database Migration
Assumption: Already cloned the repo in `$GOPATH/src` folder, and run `dep ensure`

### Creating Migration Files
#### Migration Filename Format
The ordering and direction of the migration files is determined by the filenames
used for them.  `migrate` expects the filenames of migrations to have the format:

    {version}_{title}.{extension}

The `title` of each migration is unused, and is only for readability.  Similarly,
the `extension` of the migration files is not checked by the library, and should
be an appropriate format for the database in use (`.sql` for SQL variants, for
instance).

Versions of migrations may be represented as any 64 bit unsigned integer.
All migrations are applied upward in order of increasing version number, and
downward by decreasing version number.

    {YYYYMMDDHHMMSS}_{title}.up.{extension}

    example: 20180604043223_create_con_newsfeed_components.sql

#### Writing Migrations
Migrations are defined in SQL files, which contain a set of SQL statements. Special comments are used to distinguish up and down migrations.
```sql
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE people (id int);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
```
You can put multiple statements in each block, as long as you end them with a semicolon (;).

You can alternatively set up a separator string that matches an entire line by setting sqlparse.LineSeparator. This can be used to imitate, for example, MS SQL Query Analyzer functionality where commands can be separated by a line with contents of GO. If sqlparse.LineSeparator is matched, it will not be included in the resulting migration scripts.

If you have complex statements which contain semicolons, use StatementBegin and StatementEnd to indicate boundaries:
```sql
-- +migrate Up
CREATE TABLE people (id int);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION do_something()
returns void AS $$
DECLARE
  create_query text;
BEGIN
  -- Do something here
END;
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
DROP FUNCTION do_something();
DROP TABLE people;
```
Normally each migration is run within a transaction in order to guarantee that it is fully atomic. However some SQL commands (for example creating an index concurrently in PostgreSQL) cannot be executed inside a transaction. In order to execute such a command in a migration, the migration can be run using the notransaction option:
```sql
-- +migrate Up notransaction
CREATE UNIQUE INDEX people_unique_id_idx CONCURRENTLY ON people (id);

-- +migrate Down
DROP INDEX people_unique_id_idx;
```

### Migrate Up
By default, `make migrate` will run the migration up until the latest version. However, you could add some step flag to determine the number of steps the migration should run.
```shell
$ make migrate DIRECTION=up STEP=2
```
Note: if step is provided, direction is required.

### Migrate Down
Specify the direction and migration step to rollback.
```shell
$ make migrate DIRECTION=down STEP=1
```

## Running the service
1. Make sure you already clone this repository in your `$GOPATH/src` folder.
2. Install the dependencies by running `dep ensure` in your terminal.
3. Compile the protobuf for the server example using `make proto`.
4. Create a configuration file called `config.yml` in the root directory.
5. You can copy the configuration in `config.yml.example` to your `config.yml`, or you also can change the config.
6. Run the service using `make run`.

