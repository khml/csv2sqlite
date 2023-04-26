package libc2s

import (
	"database/sql"
	"io"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase(dbFilename string) (*sql.DB, error) {
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

func CreateTable(db *sql.DB, tableName string, headerRow []string) error {
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

func InsertRecords(db *sql.DB, tableName string, data *CsvData) (int, error) {
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
		for range record {
			query += "?, "
		}
		query = strings.TrimSuffix(query, ", ") + ")"

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