package database

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type DBI interface {
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecContext(ctx context.Context, query interface{}, params ...interface{}) (orm.Result, error)

	Query(model, query interface{}, params ...interface{}) (orm.Result, error)
	QueryContext(ctx context.Context, model, query interface{}, params ...interface{}) (orm.Result, error)

	Model(model ...interface{}) *orm.Query
	ModelContext(ctx context.Context, model ...interface{}) *orm.Query

	RunInTransaction(ctx context.Context, fn func(tx *pg.Tx) error) error
}
