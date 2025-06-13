# clickhouse-migrator

<img align="right" width="125" src="assets/goose_logo.png">


Clickhouse migrator is a clickhouse database migration tool. Both a CLI and a library.

Manage your **database schema** by creating incremental SQL changes or Go functions.

#### Features

- Works for distributed clickhouse with replicated tables
- Out-of-order migrations.
- Seeding data.
- Environment variable substitution in SQL migrations.
- ... and more.

# Install

```shell
go install github.com/martirosharutyunyan/clickhouse-migrator/cmd/clickhouse-migrator@latest
```

This will install the `clickhouse-migrator` binary to your `$GOPATH/bin` directory.

# Usage

This tool is for applying migrations to distributed nodes in cluster of replicated tables in clickhouse
Which lacks the goose 

```
Clickhouse distributed migrator for distributing goose_db_version_table in cluster.
Redistributing it in case of adding or deleting nodes

Usage:
  clickhouse-migrator [flags]
  clickhouse-migrator [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      create a new migration
  down        down one migration
  down-to     down to specified version
  fix         fix migrations
  help        Help about any command
  rerun       it resets and runs migrations
  reset       reset all migrations
  status      shows the status of a migrations
  up          applying migrations to cluster
  up-by-one   up migrations one by one
  up-to       up migrations to specified version

Flags:
      --allow-missing           applies missing (out-of-order) migrations
      --cluster string          cluster name
      --db string               database name
      --dir string              migrations directory (default ".")
      --dsn string              database connection string
  -h, --help                    help for clickhouse-migrator
      --migration-type string   migration type. supported sql, go (default "sql")
      --name string             migration name
      --no-color                disable color output (NO_COLOR env variable supported)
      --no-versioning string    apply migration commands with no versioning, in file order, from directory pointed to
      --s                       use sequential numbering for new migrations
      --table string            table name for migrations versioning (default "goose_db_version")
      --timeout duration        maximum allowed duration for queries to run
      --v string                enable verbose mode

```

## License

Licensed under [MIT License](./LICENSE)