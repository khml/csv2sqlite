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

func ReadCsvFile(filename string) (*CsvData, error) {
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
