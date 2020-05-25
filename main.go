package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
	"log"
	"os"
)

type Result struct {
	Title string
	Author string
	Eval int
	TitleKana string
	AuthorKana string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loding .env file")
	}
	DBUser := os.Getenv("DB_USER_NAME")
	DBPass := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_DBNAME")
	dataSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require",
		DBHost, DBPort, DBUser, DBName, DBPass)
	db, err := sql.Open("postgres", dataSource)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(
`SELECT
a.title,
b.name AS author,
COALESCE(case when a.evaluation <> '' then a.evaluation else NULL END, '0') AS eval,
a.title_kana,
b.name_kana AS author_kana
FROM
records a
LEFT OUTER JOIN authors b ON a.author = b.id
ORDER BY
a.id ASC
`)
	if err != nil {
		log.Fatal(err)
	}
	var results []Result
	for rows.Next() {
		var result Result
		rows.Scan(&result.Title, &result.Author, &result.Eval, &result.TitleKana, &result.AuthorKana)
		results = append(results, result)
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("data")
	if err != nil {
		log.Fatal(err)
	}
	for i, data := range results {
		sheet.Cell(i, 0).Value = data.Author
		sheet.Cell(i, 1).Value = data.Title
		sheet.Cell(i, 2).SetInt(data.Eval)
		sheet.Cell(i, 3).Value = data.AuthorKana
		sheet.Cell(i, 4).Value = data.TitleKana
	}
	err = file.Save("MyArchiveBackup.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("backup output complete")
}
