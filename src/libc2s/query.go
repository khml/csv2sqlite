package libc2s

import (
	"fmt"
	"log"
	"strings"
)

// generate creat TBL query
func buildCreateTableQuery(tableName string, columnNames []string) string {
	// column definition; all column types are TEXT
	colDef := strings.Join(columnNames, " TEXT, ")

	// build create table query
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s )", tableName, colDef)

	log.Println(query)

	return query
}

// generate insert record query
func buildInsertRecordQuery(tableName string, columnNames []string) string {
	columns := strings.Join(columnNames, ",")

	colNum := len(columnNames)
	placeHolder := strings.Repeat("?, ", colNum)
	placeHolder = strings.TrimSuffix(placeHolder, ", ")

	return fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ( %s )", tableName, columns, placeHolder)
}
