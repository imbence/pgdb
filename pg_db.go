package pg_db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"strings"
)

var (
	DB  *bun.DB
	Ctx context.Context
)

func ConnectToDb(dbUrl string) error {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	DB = bun.NewDB(sqldb, pgdialect.New())
	if err := DB.Ping(); err != nil {
		return err
	}
	DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	Ctx = context.Background()
	return nil
}

func ToDb[T any](data []T, table string, schema string) (int64, error) {
	modelTableExpr := fmt.Sprintf("%s.%s", schema, table)

	onConflictString, primaryKeyString, err := getTableColumns(table, schema)
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

func getTableColumns(table, schema string) (string, string, error) {
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
