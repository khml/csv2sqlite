package libc2s

import "database/sql"

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
