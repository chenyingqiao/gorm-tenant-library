# 工具数据库分库工具包

## 安装使用

```shell
go get github.com/chenyingqiao/gorm-tenant-library
```

## 示例demo

https://git.myscrm.cn/tools/gorm-driver-test

## 初始化DB

dsn添加`dbholder`字段，用于标记当前这个db对象使用的数据库占位符。

```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=True&loc=Local&dbholder=$tenantDB",
		"root",
		"123@xxxx",
		"localhost",
		"13242",
		"",
		"utf8")
db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
```

## gin 中间件

设置后数据库名会通过这个规则去寻找
数据库名：[ServiceName]_[DbName]_[TenantCode]

```go
//数据库对应的链接配置
optItem := dto.ShardOption{
    // 服务名称
    ServiceName: "xxxx", //数据库名：[ServiceName]_[DbName]_[TenantCode]
    //库名称
    DbName:      "xxxx", //数据库名：[ServiceName]_[DbName]_[TenantCode]
    //分库标记
    TenantCode: &dto.TenantCode{
        Identify:     "oid",
        DataPosition: middleware.TenantCodeDataPositionQuery, //默认从query获取
    },
    //需要进行分库注册的db列表
    DB: db,
}

options := dto.ShardOptions{
    // 分库的方式：
    //      占位符方式 middleware.ShardDBPlaceholderWay
    //      mysql中间件方式 middleware.ShardDBKingshardWay
    //      即使用占位符也使用中间件 middleware.ShardDBPlaceholderKingshardWay
    Way:         middleware.ShardDBPlaceholderWay, 
    Options: []dto.ShardOption{
        optItem,
    },
}

engine.Use(middleware.ShardDB(options))
```

## 非中间件的注册函数

```go
//注册db
ctxDB := tgrom.WithRegisterDBContext(context.Background(), options, "1")
//$tenantDB 是dsn中设置对应的占位符
tgrom.DB(ctxDB, "$tenantDB").Model(&model.User{}).
    Where("id = ?", 1).
    Preload("ARecord").
    First(entity)
```

## model 定义上的修改

定义的model tablename前面要加上dsn上定位的占位符

```go
package model

type User struct {
	Id                 int64  `gorm:"column:id;"`
    ....
	ARecord A `gorm:"foreignkey:id;foreignkey:id"`
}

func (u *User) TableName() string {
	return "$tenantDB.user"
}
```

## dbresolver 包支持

一个db链接支持多个占位符

```go
//初始化时dbholder支持填写多个，第一个为主要的dbholder。其他的为附带的dbholder
dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=True&loc=Local&dbholder=$tenantDB",
		"root",
		"123@xxxx",
		"localhost",
		"13242",
		"",
		"utf8")
gorm.Open(driver.Open(dsn), &gorm.Config{})

//gorm分库插件注册也需要使用自定义driver
dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=True&loc=Local&dbholder=$tenantDB",
		"root",
		"123@xxxx",
		"localhost",
		"13242",
		"",
		"utf8")
db.Use(dbresolver.Register(dbresolver.Config{
    Sources: []gorm.Dialector{driver.Open(dsn)},
},
    model.XXX{}, //model中应该tablename应该是$tenantDB2.xxx
))

dbOptions := dto.ShardOption{
    ServiceName: "mars",
    DbName:      "sun-auth",
    TenantCode: &dto.TenantCode{
        Identify:     "oid",
        DataPosition: middleware.TenantCodeDataPositionQuery,
    },
    DB: db,
    OtherPlaceHolder: map[string]string{
        "$tenantDB2": "sun-auth2", //这里填写额外的数据库占位符
    },
}

options := dto.ShardOptions{
    Way: middleware.ShardDBPlaceholderWay,
    Options: []dto.ShardOption{
        dbOptions,
    },
}

engine.Use(middleware.ShardDB(options))
```

## 查询使用

查询使用有两种方式
- 自己做WithContext
- 通过包提供的获取db链接的方法获取DB对象

### 自己做WithContext

```go
//ctx 应该是gin的context中获取的 ctx.Request.Context()
db.WithContext(ctx).Model(&model.User{}).
    Where("id = ?", 1).
    Preload("ARecord").
    First(entity)
```

### 通过包提供的获取db链接的方法获取DB对象

```go
import (
    tgrom "github.com/chenyingqiao/gorm-tenant-library/gorm"
)

//$tenantDB 是dsn中设置对应的占位符
tgrom.DB(ctx.Request.Context(), "$tenantDB").Model(&model.User{}).
    Where("id = ?", 1).
    Preload("ARecord").
    First(entity)
```
