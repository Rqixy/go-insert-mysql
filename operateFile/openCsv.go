package operateFile

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// OpenCsv CSVファイルを開く関数
func OpenCsv(csvFile string) ([][]string, error) {
	// ファイルを開く
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("open error")
		return nil, err
	}

	// 最後にファイルを閉じる
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// CSVで読み込む
	r := csv.NewReader(file)

	// [][]stringでCSVファイルを受け取る
	rows, err := r.ReadAll()
	return rows, nil
}
