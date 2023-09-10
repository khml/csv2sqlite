package libc2s

import (
	"fmt"
	"log"
)

func Csv2sqlite(csvFilePath, dbFilePath, tableName string) {
	// read csv
	println("read csv ...")
	csvData, err := ReadCsv(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// connect to SQLite DB
	println("connect DB ...")
	db, err := ConnectRepository(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create TBL
	println("create DB table ...")
	err = db.CreateTbl(tableName, csvData)
	if err != nil {
		log.Fatal(err)
	}

	// insert records from CSV
	println("insert records to db ...")
	numRecords, err := db.InsertRecords(tableName, csvData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %d records into %s\n", numRecords, tableName)
}
