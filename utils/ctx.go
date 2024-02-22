package utils

import (
	"context"
	"errors"

	"github.com/chenyingqiao/gorm-tenant-library/gorm/dto"
	"gorm.io/gorm"
)

func GetDatabaseName(ctx context.Context, dbholder string) (string, error) {
	dbholder = RemovePoint(dbholder)
	dbholder += "."
	ctxDB := ctx.Value(dto.CtxKey(dbholder))
	if v, ok := ctxDB.(dto.CtxDB); ok {
		return v.Database, nil
	}
	return "", errors.New("db name is not found")
}

func GetDB(ctx context.Context, dbholder string) (*gorm.DB, error) {
	dbholder = RemovePoint(dbholder)
	dbholder += "."
	ctxDB := ctx.Value(dto.CtxKey(dbholder))
	if v, ok := ctxDB.(dto.CtxDB); ok {
		return v.DB, nil
	}
	return nil, errors.New("gorm db instance is not found")
}

func GetDbCtx(ctx context.Context, dbholder string) (dto.CtxDB, error) {
	dbholder = RemovePoint(dbholder)
	dbholder += "."
	ctxDB := ctx.Value(dto.CtxKey(dbholder))
	if v, ok := ctxDB.(dto.CtxDB); ok {
		return v, nil
	}
	return dto.CtxDB{}, errors.New("ctx db is not found")
}
