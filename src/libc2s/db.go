package libc2s

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase(dbFilename string) (*sql.DB, error) {
	// connect to SQLite DB
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func CreateTable(db *sql.DB, tableName string, headerRow []string) error {
	// generate creat TBL SQL
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for _, columnName := range headerRow {
		query += columnName + " TEXT, "
	}
	query = strings.TrimSuffix(query, ", ") + ")"
	log.Println(query)

	// exec SQL
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func InsertRecords(db *sql.DB, tableName string, data *CsvData) (int, error) {
	// counter for record num
	numRecords := 0

	// read records, one line at a time.
	for {
		record, err := data.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return numRecords, err
		}

		// generate insert record SQL
		query := "INSERT INTO " + tableName + " VALUES ("
		for range record {
			query += "?, "
		}
		query = strings.TrimSuffix(query, ", ") + ")"

		// set params
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}

		// exec insert record SQL
		_, err = db.Exec(query, args...)
		if err != nil {
			return numRecords, err
		}

		numRecords++
	}

	return numRecords, nil
}

func Csv2sqlite(csvFilePath, dbFilePath, tableName string) {
	// read csv
	println("read csv ...")
	csvData, err := ReadCsv(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// connect to SQLite DB
	println("connect DB ...")
	db, err := ConnectDatabase(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create TBL
	println("create DB table ...")
	err = CreateTable(db, tableName, csvData.HeaderRow)
	if err != nil {
		log.Fatal(err)
	}

	// insert records from CSV
	println("insert records to db ...")
	numRecords, err := InsertRecords(db, tableName, csvData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %d records into %s\n", numRecords, tableName)
}
