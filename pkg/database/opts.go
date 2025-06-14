package database

type Option func(store *ShardedClickHouseStore)

func WithQuorum(quorum int) Option {
	return func(store *ShardedClickHouseStore) {
		store.quorum = quorum
	}
}
