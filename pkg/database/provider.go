package database

import (
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/pressly/goose/v3"
	"os"
)

func NewProvider(conf *cfg.Config, opts ...goose.ProviderOption) (*goose.Provider, error) {
	store, err := NewStore(conf.DB, conf.Cluster, conf.DBName, conf.TableName)
	if err != nil {
		return nil, err
	}

	migrationsFs := os.DirFS(conf.Dir)

	opts = append(opts, goose.WithStore(store))

	return goose.NewProvider("", conf.DB, migrationsFs, opts...)
}
