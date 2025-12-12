package app

import (
	"context"
	"db_migrate_server/internal/config"
	"db_migrate_server/internal/mapping"
	"db_migrate_server/internal/pivot"
	"db_migrate_server/internal/util"
)

// LoadPlan memuat mapping root dan membuat planner.
func LoadPlan(cfg *config.Config) (*mapping.Planner, error) {
	util.Info.Println("app: loading mapping root")
	root, err := mapping.Load(*cfg)
	if err != nil {
		return nil, err
	}
	util.Info.Printf("app: mapping loaded entities=%d", len(root.Entities))

	// create planner
	return mapping.NewPlanner(root), nil
}

// InitPivot membuka repo pivot dan memastikan schema.
func InitPivot(ctx context.Context, dsn string) (*pivot.Repo, error) {
	util.Info.Println("app: connecting pivot repository")
	pivotRepo, err := pivot.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	defaultSchemaSQL := pivot.DefaultSchemaSQL()

	util.Info.Println("app: ensuring pivot schema")
	if err := pivotRepo.EnsureSchema(ctx, defaultSchemaSQL); err != nil {
		pivotRepo.Close()
		return nil, err
	}
	util.Info.Println("app: pivot schema ensured")
	return pivotRepo, nil
}
