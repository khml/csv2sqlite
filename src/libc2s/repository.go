package libc2s

import (
	"database/sql"
	"io"
)

func ConnectRepository(filepath string) (*Repository, error) {
	db, err := ConnectDatabase(filepath)
	if err != nil {
		return nil, err
	}

	return &Repository{db}, err
}

type Repository struct {
	db *sql.DB
}

func (r Repository) Close() error {
	return r.db.Close()
}

func (r Repository) CreateTbl(tblName string, csv *CsvData) error {
	query := buildCreateTableQuery(tblName, csv.HeaderRow)
	_, err := r.Exec(query)
	return err
}

func (r Repository) Exec(query string) (sql.Result, error) {
	return r.db.Exec(query)
}

func (r Repository) InsertRecords(tableName string, data *CsvData) (int, error) {
	// counter for record num
	numRecords := 0

	// build query
	query := buildInsertRecordQuery(tableName, data.HeaderRow)

	// read records, one line at a time.
	for {
		record, err := data.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return numRecords, err
		}

		// insert record
		// set params
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = v
		}

		// exec insert record SQL
		_, err = r.db.Exec(query, args...)
		if err != nil {
			return numRecords, err
		}

		numRecords++
	}

	return numRecords, nil
}
