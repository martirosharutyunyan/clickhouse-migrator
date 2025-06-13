package database

import (
	"context"
	"github.com/pressly/goose/v3"
)

func Reset(ctx context.Context, provider *goose.Provider) ([]*goose.MigrationResult, error) {
	statuses, err := provider.Status(ctx)
	if err != nil {
		return nil, err
	}
	if len(statuses) == 0 {
		return nil, nil
	}
	firstStatus := statuses[0]
	return provider.DownTo(ctx, firstStatus.Source.Version)
}
