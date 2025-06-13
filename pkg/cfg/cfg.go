package cfg

import (
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	GOOSEMIGRATIONDIR = envOr("GOOSE_MIGRATION_DIR", DefaultMigrationDir)
)

var (
	DefaultMigrationDir = "."
)

// envOr returns os.Getenv(key) if set, or else default.
func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

type Config struct {
	Dir           string
	MigrationName string
	TableName     string
	AllowMissing  bool
	NoVersioning  bool
	Sequential    bool
	NoColor       bool
	Timeout       time.Duration
	Verbose       bool
	MigrationType string
	Dsn           string
	Cluster       string
	DBName        string
	DB            *sql.DB
}

func NewConfig(cmd *cobra.Command) (*Config, error) {
	duration, err := time.ParseDuration(cmd.Flag("timeout").Value.String())
	if err != nil {
		return nil, err
	}
	dsn := cmd.Flag("dsn").Value.String()
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Config{
		Dir:           cmd.Flag("dir").Value.String(),
		TableName:     cmd.Flag("table").Value.String(),
		MigrationName: cmd.Flag("name").Value.String(),
		AllowMissing:  cmd.Flag("allow-missing").Changed,
		NoVersioning:  cmd.Flag("no-versioning").Changed,
		Sequential:    cmd.Flag("s").Changed,
		NoColor:       cmd.Flag("no-color").Changed,
		Timeout:       duration,
		Verbose:       cmd.Flag("v").Changed,
		MigrationType: cmd.Flag("migration-type").Value.String(),
		Dsn:           dsn,
		DB:            db,
		Cluster:       cmd.Flag("cluster").Value.String(),
		DBName:        cmd.Flag("db").Value.String(),
	}, nil
}
