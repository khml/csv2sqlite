package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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
	csvData, err := readCsvFile(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// SQLiteデータベースに接続する
	println("connect DB ...")
	db, err := connectDatabase(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブルを作成する
	println("create DB table ...")
	err = createTable(db, tableName, csvData.HeaderRow)
	if err != nil {
		log.Fatal(err)
	}

	// CSVファイルからレコードを挿入する
	println("insert records to db ...")
	numRecords, err := insertRecords(db, tableName, csvData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %d records into %s\n", numRecords, tableName)
}

type CsvData struct {
	HeaderRow []string
	Reader    *csv.Reader
}

func readCsvFile(filename string) (*CsvData, error) {
	// CSVファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// CSVリーダーを作成する
	reader := csv.NewReader(file)

	// ヘッダー行を読み込む
	headerRow, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("CSVファイルが空です")
		}
		return nil, err
	}

	// ヘッダー行を設定する
	reader.FieldsPerRecord = len(headerRow)
	reader.ReuseRecord = true

	return &CsvData{headerRow, reader}, nil
}

func connectDatabase(dbFilename string) (*sql.DB, error) {
	// SQLiteデータベースに接続する
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
	}

	// 接続を確認する
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB, tableName string, headerRow []string) error {
	// SQL文を作成する
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for _, columnName := range headerRow {
		query += columnName + " TEXT, "
	}
	query = strings.TrimSuffix(query, ", ") + ")"
	log.Println(query)

	// SQL文を実行する
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func insertRecords(db *sql.DB, tableName string, data *CsvData) (int, error) {
	// 挿入するレコード数を初期化する
	numRecords := 0

	// レコードを1行ずつ読み込む

	for {
		record, err := data.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return numRecords, err
		}

		// SQL文を作成する
		query := "INSERT INTO " + tableName + " VALUES ("
		for _ = range record {
			query += "?, "
		}
		query = strings.TrimSuffix(query, ", ") + ")"
		log.Println(query)

		// パラメータを設定する
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}

		// SQL文を実行する
		_, err = db.Exec(query, args...)
		if err != nil {
			return numRecords, err
		}

		numRecords++
	}

	return numRecords, nil
}
