package callback

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chenyingqiao/gorm-tenant-library/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PlaceHolder(db *gorm.DB) {
	placeHolder := db.Dialector.Name()
	ctxDB, err := utils.GetDbCtx(db.Statement.Context, placeHolder)
	if err != nil {
		return
	}

	databaseName := ctxDB.Database
	otherPlaceHolder := ctxDB.OtherPlaceHolder

	replaceHolder(db, databaseName, placeHolder, placeHolder)
	for key, item := range otherPlaceHolder {
		replaceHolder(db, item, placeHolder, key)
	}
}

func replaceHolder(db *gorm.DB, databaseName, mainPlaceHolder, placeHolder string) {
	reg, err := regexp.Compile(fmt.Sprintf(`.?\%s.?\.(.+)[\s\.]?`, utils.RemovePoint(placeHolder)))
	if err != nil {
		panic("regex error!")
	}
	//适配model
	if strings.HasPrefix(db.Statement.Table, placeHolder) {
		//需要设置
		db.Statement.Table = strings.ReplaceAll(db.Statement.Table, placeHolder, mainPlaceHolder)
		db.Statement.Table = databaseName + db.Statement.Table
	}
	if db.Statement.TableExpr != nil {
		sql := db.Statement.TableExpr.SQL
		sql = reg.ReplaceAllStringFunc(sql, func(s string) string {
			s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", databaseName))
			s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", databaseName))
			return s
		})
		db.Statement.TableExpr.SQL = sql
	}
	//适配raw
	rawSql := db.Statement.SQL.String()
	if strings.ContainsAny(rawSql, placeHolder) {
		rawSql = reg.ReplaceAllStringFunc(rawSql, func(s string) string {
			s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", databaseName))
			s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", databaseName))
			return s
		})
		db.Statement.SQL.Reset()
		db.Statement.SQL.WriteString(rawSql)
	}
	//适配joins
	for key, item := range db.Statement.Joins {
		name := reg.ReplaceAllStringFunc(item.Name, func(s string) string {
			s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", databaseName))
			s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", databaseName))
			return s
		})
		item.Name = name
		db.Statement.Joins[key] = item
	}
	//适配条件
	for key, item := range db.Statement.Clauses {
		if wheres, ok := item.Expression.(clause.Where); ok {
			for wk, wi := range wheres.Exprs {
				wi := expressionProcess(reg, wi, databaseName, placeHolder)
				wheres.Exprs[wk] = wi
			}
			item.Expression = wheres
			db.Statement.Clauses[key] = item
			continue
		}
	}
}

func expressionProcess(reg *regexp.Regexp, item clause.Expression, database string, placeHolder string) clause.Expression {
	if expr, ok := item.(clause.Expr); ok {
		sql := reg.ReplaceAllStringFunc(expr.SQL, func(s string) string {
			s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", database))
			s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", database))
			return s
		})
		expr.SQL = sql
		item = expr
	}
	if expr, ok := item.(clause.NamedExpr); ok {
		sql := reg.ReplaceAllStringFunc(expr.SQL, func(s string) string {
			s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", database))
			s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", database))
			return s
		})
		expr.SQL = sql
		item = expr
	}
	if expr, ok := item.(clause.IN); ok {
		if colStr, isStr := expr.Column.(string); isStr {
			sql := reg.ReplaceAllStringFunc(colStr, func(s string) string {
				s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", database))
				s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", database))
				return s
			})
			expr.Column = sql
			item = expr
		}
	}
	if expr, ok := item.(clause.Eq); ok {
		if colStr, isStr := expr.Column.(string); isStr {
			sql := reg.ReplaceAllStringFunc(colStr, func(s string) string {
				s = strings.ReplaceAll(s, placeHolder, fmt.Sprintf("`%s`.", database))
				s = strings.ReplaceAll(s, fmt.Sprintf("`%s`.", utils.RemovePoint(placeHolder)), fmt.Sprintf("`%s`.", database))
				return s
			})
			expr.Column = sql
			item = expr
		}
	}
	return item
}
