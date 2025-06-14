-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ${CH_DB}.game_counts_local ON CLUSTER '${CH_CLUSTER}'
(
    user_id               UUID,
    game_counts           UInt32
)
ENGINE = ReplicatedSummingMergeTree('/clickhouse/${CH_DB}/tables/{shard}/game_counts', '{replica}')
ORDER BY user_id
SETTINGS index_granularity = 8192;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ${CH_DB}.game_counts_local ON CLUSTER '${CH_CLUSTER}' SYNC;
-- +goose StatementEnd
