package trace

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	gormTraceStartTime = "gorm-trace-start-time"
	gormTraceContext   = "gorm-trace-context"
)

func WithContext(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.Set(gormTraceContext, ctx)
}

func SetLogger(db *gorm.DB, logger Logger) {

	after := afterFunc(logger)
	cb := db.Callback()
	cb.Create().Before("gorm:before_create").Register("gorm-trace:before_create", before)
	cb.Create().After("gorm:after_create").Register("gorm-trace:after_create", after)
	cb.Update().Before("gorm:before_update").Register("gorm-trace:before_update", before)
	cb.Update().After("gorm:after_update").Register("gorm-trace:after_update", after)
	cb.Delete().Before("gorm:before_delete").Register("gorm-trace:before_delete", before)
	cb.Delete().After("gorm:after_delete").Register("gorm-trace:after_delete", after)
	cb.Query().Before("gorm:query").Register("gorm-trace:before_query", before)
	cb.Query().After("gorm:after_query").Register("gorm-trace:after_query", after)
	cb.RowQuery().Before("gorm:row_query").Register("gorm-trace:before_row_query", before)
	cb.RowQuery().After("gorm:row_query").Register("gorm-trace:after_row_query", after)
}

func before(scope *gorm.Scope) {
	scope.Set(gormTraceStartTime, time.Now())
}

func afterFunc(logger Logger) func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		var ctx context.Context
		if ictx, ok := scope.Get(gormTraceContext); ok {
			ctx = ictx.(context.Context)
		} else {
			ctx = context.Background()
		}

		if scope.DB().Error != nil {
			logger.Print(ctx, scope.DB().Error)
		}

		startTime, _ := scope.Get(gormTraceStartTime)
		t := time.Now().Sub(startTime.(time.Time))
		logger.Print(ctx, "sql", t, scope.SQL, scope.SQLVars, scope.DB().RowsAffected)
	}
}
