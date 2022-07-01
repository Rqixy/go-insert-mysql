package operateFile

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func OpenText(file string) string {
	fmt.Println("ファイルを開きます。")

	var f, err = os.Open("text.txt")
	if err != nil {
		fmt.Println("error")
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	b, err := ioutil.ReadAll(f)

	return string(b)
}
