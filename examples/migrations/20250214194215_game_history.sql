-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ${CH_DB}.game_history_local ON CLUSTER '${CH_CLUSTER}'
(
    id            UUID DEFAULT generateUUIDv4(),
    created_at    DATETIME,
    user_id       UUID,
    car_id        UUID,
    score         UInt32
)
ENGINE = ReplicatedMergeTree('/clickhouse/${CH_DB}/tables/{shard}/game_history', '{replica}')
    PARTITION BY toYYYYMMDD(created_at)
    ORDER BY score
SETTINGS index_granularity = 8192;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ${CH_DB}.game_history_local ON CLUSTER '${CH_CLUSTER}' SYNC;
-- +goose StatementEnd
