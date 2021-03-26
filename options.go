package workflow

import (
	"context"

	gormV2 "gorm.io/gorm"
)

const (
	GormDBKey = "gorm-db"
)

type Options interface {
	Apply(data *WorkData)
}

type applyFunc func(data *WorkData)

func (f applyFunc) Apply(data *WorkData) {
	f(data)
}

func MustGetGormTx(data *WorkData) *gormV2.DB {
	return data.MustGet(GormDBKey).(*gormV2.DB)
}

func WithGormV2(db *gormV2.DB) Options {
	return applyFunc(func(data *WorkData) {
		data.workBegin.Add(func(ctx context.Context, data *WorkData) error {
			tx := db.Begin()
			if tx.Error != nil {
				return tx.Error
			}
			data.Set(GormDBKey, tx)
			return nil
		})
		data.workCommit.Add(func(ctx context.Context, data *WorkData) error {
			return MustGetGormTx(data).Commit().Error
		})
		data.workRollback.Add(func(ctx context.Context, data *WorkData) error {
			return MustGetGormTx(data).Rollback().Error
		})
	})
}

func WithBegin(f Event) Options {
	return applyFunc(func(data *WorkData) {
		data.workBegin.Add(f)
	})
}

func WithBeforeCommit(f Event) Options {
	return applyFunc(func(data *WorkData) {
		data.workBeforeCommit.Add(f)
	})
}

func WithCommit(f Event) Options {
	return applyFunc(func(data *WorkData) {
		data.workCommit.Add(f)
	})
}

func WithRollback(f Event) Options {
	return applyFunc(func(data *WorkData) {
		data.workRollback.Add(f)
	})
}

func WithFinish(f Event) Options {
	return applyFunc(func(data *WorkData) {
		data.workFinish.Add(f)
	})
}

func WithInterface(param interface{}) Options {
	return applyFunc(func(data *WorkData) {
		if f, ok := param.(interface {
			WorkBegin(context.Context, *WorkData) error
		}); ok {
			WithBegin(f.WorkBegin).Apply(data)
		}

		if f, ok := param.(interface {
			WorkBeforeCommit(context.Context, *WorkData) error
		}); ok {
			WithBeforeCommit(f.WorkBeforeCommit).Apply(data)
		}

		if f, ok := param.(interface {
			WorkCommit(context.Context, *WorkData) error
		}); ok {
			WithCommit(f.WorkCommit).Apply(data)
		}

		if f, ok := param.(interface {
			WorkRollback(context.Context, *WorkData) error
		}); ok {
			WithRollback(f.WorkRollback).Apply(data)
		}

		if f, ok := param.(interface {
			WorkFinish(context.Context, *WorkData) error
		}); ok {
			WithFinish(f.WorkFinish).Apply(data)
		}
	})
}
