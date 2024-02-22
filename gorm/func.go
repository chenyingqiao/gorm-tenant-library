package gorm

import (
	"context"
	"fmt"

	"github.com/chenyingqiao/gorm-tenant-library/gorm/callback"
	"github.com/chenyingqiao/gorm-tenant-library/gorm/dto"
	"github.com/chenyingqiao/gorm-tenant-library/utils"
	"gorm.io/gorm"
)

// 获取数据库实例
func DB(ctx context.Context, dbholder string) *gorm.DB {
	db, err := utils.GetDB(ctx, dbholder)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func WithRegisterDBContext(ctx context.Context, options dto.ShardOptions, tenantCode string) context.Context {
	for _, item := range options.Options {
		dbName := fmt.Sprintf("%s_%s_%s", item.ServiceName, item.DbName, tenantCode)
		//设置数据库callback
		if item.DB.Callback().Query().Get("database_query") == nil {
			item.DB.Callback().Query().Before("*").Register("database_query", callback.PlaceHolder)
			item.DB.Callback().Update().Before("*").Register("database_update", callback.PlaceHolder)
			item.DB.Callback().Delete().Before("*").Register("database_delete", callback.PlaceHolder)
			item.DB.Callback().Create().Before("*").Register("database_create", callback.PlaceHolder)
			item.DB.Callback().Raw().Before("*").Register("database_raw", callback.PlaceHolder)
			item.DB.Callback().Row().Before("*").Register("database_row", callback.PlaceHolder)
		}

		//获取额外的替换占位符
		otherPh, err := GetOtherPlaceHolderDbName(item, tenantCode)
		if err != nil {
			continue
		}

		//预先放置数据库相关信息
		ctxDbWith := context.WithValue(ctx, dto.CtxKey(item.DB.Dialector.Name()), dto.CtxDB{
			Database:         dbName,
			OtherPlaceHolder: otherPh,
		})

		//设置请求独立的db
		newDb := item.DB.WithContext(ctxDbWith)

		//设置需要查询的数据库名称
		ctx = context.WithValue(ctx, dto.CtxKey(item.DB.Dialector.Name()), dto.CtxDB{
			Database:         dbName,
			DB:               newDb,
			OtherPlaceHolder: otherPh,
		})
	}
	return ctx
}

func GetOtherPlaceHolderDbName(options dto.ShardOption, tenantCode string) (map[string]string, error) {
	result := map[string]string{}

	otherPlaceHolder := options.OtherPlaceHolder

	for key, item := range otherPlaceHolder {
		result[key] = fmt.Sprintf("%s_%s_%s", options.ServiceName, item, tenantCode)
	}

	return result, nil
}
