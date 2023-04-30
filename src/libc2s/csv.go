package libc2s

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CsvData struct {
	HeaderRow []string
	Reader    *csv.Reader
}

func ReadCsv(pathToCsv string) (*CsvData, error) {
	// open csv file
	file, err := os.Open(pathToCsv)
	if err != nil {
		return nil, err
	}

	// spawn csv reader
	reader := csv.NewReader(file)

	// read only header row
	headerRow, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("CSV file (%s) is empty", pathToCsv)
		}
		return nil, err
	}

	// set column num
	reader.FieldsPerRecord = len(headerRow)
	reader.ReuseRecord = true

	return &CsvData{headerRow, reader}, nil
}
