package operateDb

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rqixy/insertMysql/operateFile"
	"log"
	"os"
)

type QuizInfo struct {
	QuizID     int
	Quiz       string
	Answer     string
	Commentary string
	AnswerID   int
}

type Answers struct {
	AnswerId   int
	Correct    string
	Incorrect1 string
	Incorrect2 string
	Incorrect3 string
}

var db *sql.DB

func DBConnect() (*sql.DB, error) {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error open .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbAddr := os.Getenv("DB_ADDRESS")
	dbName := os.Getenv("DB_NAME")

	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPass,
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("connected!!")

	return db, err
}

func InsertAnswers() (sql.Result, error) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO answers (answer_id, correct, incorrect_1, incorrect_2, incorrect_3) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()

	answerCSV, err := operateFile.OpenCsv("CSVFile/insertAnswers.csv")

	var result sql.Result

	for id, v := range answerCSV {
		answers := Answers{
			AnswerId:   id + 1,
			Correct:    v[1],
			Incorrect1: v[2],
			Incorrect2: v[3],
			Incorrect3: v[4],
		}

		result, err = stmt.Exec(
			&answers.AnswerId,
			&answers.Correct,
			&answers.Incorrect1,
			&answers.Incorrect2,
			&answers.Incorrect3,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func InsertQuestion() (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction begin error: %v", err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO quiz_info (quiz_id, quiz, answer, commentary, answer_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("statement error: %v", err)
	}
	//defer db.Close()

	questionCSV, err := operateFile.OpenCsv("CSVFile/insertQuestion.csv")
	if err != nil {
		return nil, fmt.Errorf("openCSV error: %v", err)
	}

	var result sql.Result
	ansID, err := AnswerID()
	if err != nil {
		return nil, fmt.Errorf("AnswerID error: %v", err)
	}
	var quizInfos []QuizInfo
	var quizInfo QuizInfo
	for k, v := range questionCSV {
		quizInfo = QuizInfo{
			QuizID:     k + 1,
			Quiz:       v[1],
			Answer:     v[2],
			Commentary: v[3],
			AnswerID:   ansID[k],
		}
		quizInfos = append(quizInfos, quizInfo)
	}

	for _, quiz := range quizInfos {
		result, err = stmt.Exec(
			quiz.QuizID,
			quiz.Quiz,
			quiz.Answer,
			quiz.Commentary,
			quiz.AnswerID,
		)
		if err != nil {
			return nil, fmt.Errorf("statement execute error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commit error: %v", err)
	}

	return result, nil
}

func AnswerID() ([]int, error) {
	var answersID []int

	ids, err := db.Query("SELECT answer_id FROM answers")
	if err != nil {
		log.Fatal(err)
	}

	for ids.Next() {
		var answerID int
		if err := ids.Scan(&answerID); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		answersID = append(answersID, answerID)
	}
	if err := ids.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return answersID, nil
}

func DBClose() error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("DB close error: %v", err)
	}
	return nil
}
