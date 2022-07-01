package main

import (
	"github.com/rqixy/insertMysql/operateDb"
	"log"
)

func main() {
	var err error
	_, err = operateDb.DBConnect()
	if err != nil {
		log.Fatal(err)
	}

	_, err = operateDb.InsertAnswers()
	if err != nil {
		log.Fatal(err)
	}

	_, err = operateDb.InsertQuestion()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := operateDb.DBClose()
		if err != nil {
			log.Fatal(err)
		}
	}()

}
