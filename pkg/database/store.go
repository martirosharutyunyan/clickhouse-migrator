package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/types"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3/database" // Import the goose database package
)

// Ensure that ShardedClickHouseStore implements the Store interface.
var _ database.Store = (*ShardedClickHouseStore)(nil)

// ShardedClickHouseStore is a Store implementation for ClickHouse in a sharded environment.
type ShardedClickHouseStore struct {
	tableName string
	cluster   string // ClickHouse cluster name.  Important for Distributed tables.
	db        *sql.DB
	dbName    string
}

// NewShardedClickHouseStore creates a new ShardedClickHouseStore.
// cluster: The name of the ClickHouse cluster.
func NewShardedClickHouseStore(db *sql.DB, cluster, dbName string, opts ...Option) *ShardedClickHouseStore {
	store := &ShardedClickHouseStore{
		tableName: "goose_db_version", // Default table name, same as standard goose
		cluster:   cluster,
		db:        db,
		dbName:    dbName,
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

// Tablename returns the name of the version table.
func (s *ShardedClickHouseStore) Tablename() string {
	return s.tableName
}

// CreateVersionTable creates the version table in ClickHouse.  Creates both a local table
// and a distributed table.
func (s *ShardedClickHouseStore) CreateVersionTable(ctx context.Context, db database.DBTxConn) error {
	//check if the db is *sql.DB.  If it is, use it.  If it is *sql.Tx, return error
	if _, ok := db.(*sql.Tx); ok {
		return fmt.Errorf("CreateVersionTable needs *sql.DB, not *sql.Tx")
	}

	// Create a local table on each shard.  ReplicatedMergeTree is a good choice for this.
	localTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s_local ON CLUSTER '%s'
		(
			version             Int64,
			is_applied          UInt8,
			t                   DateTime
		)
		ENGINE = ReplicatedMergeTree('/clickhouse/%s/tables/{shard}/%s', '{replica}')
		ORDER BY (version)
		SETTINGS index_granularity = 1`,
		s.dbName, s.tableName, s.cluster, s.dbName, s.tableName,
	)

	_, err := s.db.ExecContext(ctx, localTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create local version table: %w", err)
	}

	// Create a distributed table that points to the local table.
	distributedTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s ON CLUSTER '%s' AS %s_local
		ENGINE = Distributed('%s', '%s', '%s_local', cityHash64(version))`,
		s.dbName, s.tableName, s.cluster, s.tableName, s.cluster, s.dbName, s.tableName,
	)

	_, err = s.db.ExecContext(ctx, distributedTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create distributed version table: %w", err)
	}

	list, err := s.ListMigrations(ctx, db)
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return nil
	}

	// Insert the initial version (0) into the distributed table.
	_, err = s.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s.%s (version, is_applied, t) VALUES (0, 1, now())", s.dbName, s.tableName))
	if err != nil {
		return fmt.Errorf("failed to insert initial version: %w", err)
	}

	return nil
}

// Insert inserts a version id into the version table.
func (s *ShardedClickHouseStore) Insert(ctx context.Context, db database.DBTxConn, req database.InsertRequest) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (version, is_applied, t) VALUES (?, 1, now())", s.dbName, s.tableName)
	_, err := db.ExecContext(ctx, query, req.Version)
	if err != nil {
		return fmt.Errorf("failed to insert version %d: %w", req.Version, err)
	}
	return nil
}

func (s *ShardedClickHouseStore) BulkInsert(ctx context.Context, versions []int64, quorum int64) error {
	query := fmt.Sprintf(`INSERT INTO %s.%s (version, is_applied, t) 
		SETTINGS 
			insert_distributed_sync = 1,
			insert_quorum = %d, 
    		insert_quorum_parallel = 0 
 		VALUES `, s.dbName, s.tableName, quorum,
	)

	var versionsAny []any
	for i, version := range versions {
		if i == len(versions)-1 {
			query += "(?, 1, now())"
			versionsAny = append(versionsAny, version)
			continue
		}
		query += "(?, 1, now()), "
		versionsAny = append(versionsAny, version)
	}

	_, err := s.db.ExecContext(ctx, query, versionsAny...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert version %d: %w", versions, err)
	}
	return nil
}

func (s *ShardedClickHouseStore) GetMigrationsWithShards(ctx context.Context) ([]types.Migration, error) {
	// Important:  Order by version DESC.
	query := fmt.Sprintf("SELECT version, is_applied, shardNum() as shard_num FROM %s.%s ORDER BY version DESC SETTINGS select_sequential_consistency = 1", s.dbName, s.tableName)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list migrations: %w", err)
	}
	defer rows.Close()

	var migrations []types.Migration
	for rows.Next() {
		var (
			version   int64
			isApplied uint8 // ClickHouse uses UInt8
			shardNum  uint8 // ClickHouse uses UInt8
		)
		if err := rows.Scan(&version, &isApplied, &shardNum); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		migrations = append(migrations, types.Migration{
			Version:     version,
			ShardNumber: int(shardNum),
			IsApplied:   isApplied == 1,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return migrations, nil
}

func (s *ShardedClickHouseStore) GetReplicaCount(ctx context.Context) (int64, error) {
	query := fmt.Sprintf(`SELECT
		count() AS replica_count
		FROM system.replicas
		WHERE database = '%s' AND table = '%s_local';`,
		s.dbName,
		s.tableName,
	)

	row := s.db.QueryRowContext(ctx, query)
	var replicaCount int64
	if err := row.Scan(&replicaCount); err != nil {
		return 0, fmt.Errorf("failed to get replica count: %w", err)
	}

	return replicaCount, nil
}

// Delete deletes a version id from the version table.
func (s *ShardedClickHouseStore) Delete(ctx context.Context, db database.DBTxConn, version int64) error {
	query := fmt.Sprintf("ALTER TABLE %s.%s_local ON CLUSTER '%s' DELETE WHERE version = ?", s.dbName, s.tableName, s.cluster)
	_, err := db.ExecContext(ctx, query, version)
	if err != nil {
		return fmt.Errorf("failed to delete version %d: %w", version, err)
	}
	return nil
}

func (s *ShardedClickHouseStore) Truncate(ctx context.Context) error {
	query := fmt.Sprintf("TRUNCATE %s.%s_local ON CLUSTER '%s' SETTINGS alter_sync=2;", s.dbName, s.tableName, s.cluster)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate %s: %w", s.tableName, err)
	}
	return nil
}

// GetMigration retrieves a single migration by version id.
func (s *ShardedClickHouseStore) GetMigration(ctx context.Context, db database.DBTxConn, version int64) (*database.GetMigrationResult, error) {
	query := fmt.Sprintf("SELECT t, is_applied FROM %s.%s WHERE version = ?", s.dbName, s.tableName)
	row := db.QueryRowContext(ctx, query, version)

	var result database.GetMigrationResult
	var isApplied uint8 // ClickHouse uses UInt8 for boolean
	err := row.Scan(&result.Timestamp, &isApplied)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrVersionNotFound
		}
		return nil, fmt.Errorf("failed to get migration version %d: %w", version, err)
	}
	result.IsApplied = isApplied == 1
	return &result, nil
}

// GetLatestVersion retrieves the last applied migration version.
func (s *ShardedClickHouseStore) GetLatestVersion(ctx context.Context, db database.DBTxConn) (int64, error) {
	query := fmt.Sprintf("SELECT max(version) FROM %s.%s", s.dbName, s.tableName)
	row := db.QueryRowContext(ctx, query)

	var version int64
	err := row.Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, database.ErrVersionNotFound // Special case: no migrations at all.
		}
		return 0, fmt.Errorf("failed to get latest version: %w", err)
	}
	return version, nil
}

// ListMigrations retrieves all migrations sorted in descending order by id or timestamp.
func (s *ShardedClickHouseStore) ListMigrations(ctx context.Context, db database.DBTxConn) ([]*database.ListMigrationsResult, error) {
	// Important:  Order by version DESC.
	query := fmt.Sprintf("SELECT version, is_applied FROM %s.%s ORDER BY version DESC", s.dbName, s.tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list migrations: %w", err)
	}
	defer rows.Close()

	var migrations []*database.ListMigrationsResult
	for rows.Next() {
		var (
			version   int64
			isApplied uint8 // ClickHouse uses UInt8
		)
		if err := rows.Scan(&version, &isApplied); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		migrations = append(migrations, &database.ListMigrationsResult{
			Version:   version,
			IsApplied: isApplied == 1,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return migrations, nil
}

// NewStore is a factory function that returns a new Store implementation for the given dialect.
//
// It is recommended to use this function instead of directly instantiating Store implementations.
//
// The dialect parameter is a string that specifies the database dialect.  If an empty string is
// provided, it attempts to auto-detect the dialect from the database connection.
func NewStore(db *sql.DB, clusterName, dbName, tableName string, opts ...Option) (*ShardedClickHouseStore, error) {
	if clusterName == "" || dbName == "" {
		return nil, fmt.Errorf("cluster name and db name must be provided")
	}
	store := NewShardedClickHouseStore(db, clusterName, dbName, opts...)
	if tableName != "" {
		store.tableName = tableName
	}

	if err := store.CreateVersionTable(context.Background(), db); err != nil {
		return nil, err
	}

	return store, nil
}
