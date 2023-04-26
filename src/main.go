package main

import (
	"fmt"
	"log"
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

	// CSVファイルを開き、CSVファイルをパースする
	println("read csv ...")
	csvData, err := libc2s.ReadCsvFile(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// SQLiteデータベースに接続する
	println("connect DB ...")
	db, err := libc2s.ConnectDatabase(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブルを作成する
	println("create DB table ...")
	err = libc2s.CreateTable(db, tableName, csvData.HeaderRow)
	if err != nil {
		log.Fatal(err)
	}

	// CSVファイルからレコードを挿入する
	println("insert records to db ...")
	numRecords, err := libc2s.InsertRecords(db, tableName, csvData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %d records into %s\n", numRecords, tableName)
}
