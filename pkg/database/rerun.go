package database

import (
	"context"
	"github.com/pressly/goose/v3"
)

func Rerun(ctx context.Context, provider *goose.Provider) ([]*goose.MigrationResult, error) {
	_, err := Reset(ctx, provider)
	if err != nil {
		return nil, err
	}

	return provider.Up(ctx)
}
