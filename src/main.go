package main

import (
	"fmt"
	"os"

	"csv2sqlite/libc2s"
)

func main() {
	// コマンドライン引数を取得する
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <csv_file_path> <table_name> <database_file_path>\n", os.Args[0])
		os.Exit(1)
	}
	csvFilePath := os.Args[1]
	tableName := os.Args[2]
	dbFilePath := os.Args[3]

	libc2s.Csv2sqlite(csvFilePath, dbFilePath, tableName)
}
