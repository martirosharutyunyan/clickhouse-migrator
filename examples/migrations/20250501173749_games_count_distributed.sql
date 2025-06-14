-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ${CH_DB}.game_counts ON CLUSTER '${CH_CLUSTER}'
AS game_counts_local
    ENGINE = Distributed('${CH_CLUSTER}', '${CH_DB}', game_counts_local, cityHash64(auth_login_history_id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ${CH_DB}.game_counts ON CLUSTER '${CH_CLUSTER}' SYNC;
-- +goose StatementEnd
