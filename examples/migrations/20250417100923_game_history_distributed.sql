-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ${CH_DB}.game_history ON CLUSTER '${CH_CLUSTER}'
AS game_history_local
    ENGINE = Distributed('${CH_CLUSTER}', '${CH_DB}', game_history_local, cityHash64(created_at));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ${CH_DB}.game_history ON CLUSTER '${CH_CLUSTER}' SYNC;
-- +goose StatementEnd
