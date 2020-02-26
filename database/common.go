package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/PPIO/pi-cloud-monitor-backend/logger"
)

var log = logger.New("database")

var (
	ErrNoRecord      = fmt.Errorf("no record")
	ErrOpt           = fmt.Errorf("database operation error")
	ErrInternal      = fmt.Errorf("internal error")
	ErrNoRowAffected = fmt.Errorf("no row affected")
	ErrParam         = fmt.Errorf("invalid parameter")
)

// DBEnsure ensures the sql's execution has taken effect.
func DBEnsure(ctx context.Context, stmt *sql.Stmt, params ...interface{}) error {
	result, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		return fmt.Errorf("mysql: fail to execute: %v", err)
	}
	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql: fail to query the num of affected rows: %v", err)
	}
	if row <= 0 {
		return ErrNoRowAffected
	}
	return nil
}

// addFileters 用于在查询语句中添加 where 条件
func addFileters(sql string, filters map[string]string) string {
	addedSQL := ""
	if len(filters) > 0 {
		addedSQL += "WHERE "
		var filterList []string
		for key, value := range filters {
			if key == "nat_type" || key == "uid" {
				f := fmt.Sprintf(`%s = %s`, key, value)
				filterList = append(filterList, f)
			} else {
				f := fmt.Sprintf(`%s = "%s"`, key, value)
				filterList = append(filterList, f)
			}
		}
		addedSQL += strings.Join(filterList, " AND ")
	}
	nSQL := fmt.Sprintf(sql, addedSQL)
	return nSQL
}
