package database

import (
	"context"
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/types"
	"github.com/pressly/goose/v3"
	"slices"
	"strconv"
	"strings"
)

func Reshard(
	ctx context.Context,
	conf *cfg.Config,
	provider *goose.Provider,
	quorum int,
) ([]types.Migration, error) {
	store, err := NewStore(conf.DB, conf.Cluster, conf.DBName, conf.TableName, WithQuorum(quorum))
	if err != nil {
		return nil, err
	}

	migrations, err := store.GetMigrationsWithShards(ctx)
	if err != nil {
		return nil, err
	}

	sources := provider.ListSources()

	err = store.Truncate(ctx)
	if err != nil {
		return nil, err
	}

	insertedMigrations := []int64{0}
	var migrationsWithSources []types.Migration
	for _, migration := range migrations {
		if migration.IsApplied && migration.Version != 0 {
			sourceIndex := slices.IndexFunc(sources, func(source *goose.Source) bool {
				return strings.Contains(source.Path, strconv.Itoa(int(migration.Version)))
			})
			if sourceIndex == -1 {
				panic(fmt.Errorf("migration %s not found in source %s", migration.Version, sources))
			}
			insertedMigrations = append(insertedMigrations, migration.Version)
			migrationsWithSources = append(migrationsWithSources, types.Migration{
				Version:   migration.Version,
				IsApplied: migration.IsApplied,
				Source:    sources[sourceIndex].Path,
			})
		}
	}

	err = store.BulkInsert(ctx, insertedMigrations)
	if err != nil {
		return nil, err
	}

	//time.Sleep(time.Second)

	reshardedMigrations, err := store.GetMigrationsWithShards(ctx)
	if err != nil {
		return nil, err
	}

	reshardedMigrations = slices.DeleteFunc(reshardedMigrations, func(i types.Migration) bool {
		return i.Version == 0
	})

	for i, migration := range migrationsWithSources {
		reshardedMigrations[i].Source = migration.Source
	}

	slices.SortFunc(reshardedMigrations, func(a, b types.Migration) int {
		return int(a.Version - b.Version)
	})

	return reshardedMigrations, nil
}
