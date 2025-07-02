package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"os"
	"strings"
)

func ConnectToDb(dbUrl string) (error, bun.DB, context.Context) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	var DB = bun.NewDB(sqldb, pgdialect.New())
	if err := DB.Ping(); err != nil {
		return err, *DB, nil
	}
	if os.Getenv("DEBUG_SQL") == "true" {
		DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	var Ctx = context.Background()
	return nil, *DB, Ctx
}

func ToDb[T any](Ctx context.Context, DB bun.DB, data []T, table string, schema string) (int64, error) {
	modelTableExpr := fmt.Sprintf("%s.%s", schema, table)

	onConflictString, primaryKeyString, err := getTableColumns(Ctx, DB, table, schema)
	if err != nil {
		return 0, err
	}
	sqlRes, err := DB.NewInsert().
		ModelTableExpr(modelTableExpr).
		Model(&data).
		On(fmt.Sprintf("CONFLICT (%s) DO UPDATE", primaryKeyString)).
		Set(onConflictString).
		Exec(Ctx)
	if err != nil {
		return 0, err
	}
	return sqlRes.RowsAffected()
}

func getTableColumns(Ctx context.Context, DB bun.DB, table string, schema string) (string, string, error) {
	var columns []string
	var primaryKey []string

	if err := DB.NewSelect().
		ColumnExpr("column_name").
		TableExpr("information_schema.columns").
		Where("table_name = ? and table_schema = ?", table, schema).
		Scan(Ctx, &columns); err != nil {
		return "", "", err
	}

	if err := DB.NewSelect().
		ColumnExpr("column_name").
		TableExpr("information_schema.key_column_usage").
		Where("table_name = ? and table_schema = ?", table, schema).
		Scan(Ctx, &primaryKey); err != nil {
		return "", "", err
	}

	excludedColumns := make([]string, 0, len(columns))
	for _, col := range columns {
		excludedColumns = append(excludedColumns, fmt.Sprintf("%s = excluded.%s", col, col))
	}

	return strings.Join(excludedColumns, ", "), strings.Join(primaryKey, ", "), nil
}
