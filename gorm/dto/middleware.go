package dto

import "gorm.io/gorm"

type CtxKey string

type ShardOption struct {
	ServiceName      string
	DbName           string
	TenantCode       *TenantCode
	DB               *gorm.DB
	OtherPlaceHolder map[string]string //适配多个PlaceHolder
}

type ShardOptions struct {
	Way     uint8
	Options []ShardOption
}

type TenantCode struct {
	Identify     string
	DataPosition uint
}

type CtxDB struct {
	Database         string
	DB               *gorm.DB
	OtherPlaceHolder map[string]string
}
