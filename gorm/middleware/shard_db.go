package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/chenyingqiao/gorm-tenant-library/gorm/callback"
	"github.com/chenyingqiao/gorm-tenant-library/gorm/dto"
	"github.com/gin-gonic/gin"
)

const (
	// 分库方式
	NoShardDB                      = 0
	ShardDBPlaceholderWay          = 1
	ShardDBKingshardWay            = 2
	ShardDBPlaceholderKingshardWay = 3

	//公司标识数据位置
	TenantCodeDataPositionQuery       = 1
	TenantCodeDataPositionHeader      = 2
	TenantCodeDefaultIdentify         = "oid"
	TenantCodeDataDefaultDataPosition = TenantCodeDataPositionQuery
)

func ShardDB(options dto.ShardOptions) gin.HandlerFunc {
	if options.Way == ShardDBPlaceholderWay {
		return func(ctx *gin.Context) {
			for _, item := range options.Options {
				dbName, err := GetDatabaseName(ctx, item)
				if err != nil {
					continue
				}
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
				otherPh, err := GetOtherPlaceHolderDbName(ctx, item)
				if err != nil {
					continue
				}

				//预先放置数据库信息
				ctxDbWith := context.WithValue(ctx.Request.Context(), dto.CtxKey(item.DB.Dialector.Name()), dto.CtxDB{
					Database:         dbName,
					OtherPlaceHolder: otherPh,
				})

				//设置请求独立的db
				newDb := item.DB.WithContext(ctxDbWith)

				//设置需要查询的数据库名称
				ctxKV := context.WithValue(ctx.Request.Context(), dto.CtxKey(item.DB.Dialector.Name()), dto.CtxDB{
					Database:         dbName,
					DB:               newDb,
					OtherPlaceHolder: otherPh,
				})
				ctx.Request = ctx.Request.WithContext(ctxKV)
			}
			ctx.Next()
		}
	}
	return emptyWay
}

func emptyWay(ctx *gin.Context) {
	ctx.Next()
}

func GetOtherPlaceHolderDbName(ctx *gin.Context, options dto.ShardOption) (map[string]string, error) {
	result := map[string]string{}

	otherPlaceHolder := options.OtherPlaceHolder

	if options.TenantCode == nil {
		options.TenantCode = &dto.TenantCode{
			Identify:     TenantCodeDefaultIdentify,
			DataPosition: TenantCodeDataDefaultDataPosition,
		}
	}

	oid := ""
	if options.TenantCode.DataPosition == TenantCodeDataPositionQuery {
		oid = ctx.Query(options.TenantCode.Identify)
	} else if options.TenantCode.DataPosition == TenantCodeDataPositionHeader {
		oid = ctx.GetHeader(options.TenantCode.Identify)
	} else {
		return result, errors.New("tenant code is empty")
	}

	for key, item := range otherPlaceHolder {
		result[key] = fmt.Sprintf("%s_%s_%s", options.ServiceName, item, oid)
	}

	return result, nil
}

func GetDatabaseName(ctx *gin.Context, options dto.ShardOption) (string, error) {
	if options.TenantCode == nil {
		options.TenantCode = &dto.TenantCode{
			Identify:     TenantCodeDefaultIdentify,
			DataPosition: TenantCodeDataDefaultDataPosition,
		}
	}

	oid := ""
	if options.TenantCode.DataPosition == TenantCodeDataPositionQuery {
		oid = ctx.Query(options.TenantCode.Identify)
	} else if options.TenantCode.DataPosition == TenantCodeDataPositionHeader {
		oid = ctx.GetHeader(options.TenantCode.Identify)
	} else {
		return "", errors.New("tenant code is empty")
	}

	if oid == "" {
		return "", errors.New("tenant code is empty")
	}

	return fmt.Sprintf("%s_%s_%s", options.ServiceName, options.DbName, oid), nil
}
