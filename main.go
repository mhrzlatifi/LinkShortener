package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"linkShortener/DB"
	"linkShortener/Routes"
	"log"
	"net/http"
)

func main() {

	username := "root"
	password := ""
	host := "localhost"
	port := "3306"
	dbName := "link_shortener"

	DB.MYSQL, _ = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName))

	defer DB.MYSQL.Close()

	mux := Routes.SetupRoutes()
	err := http.ListenAndServe(":8080", mux)
	errCheck(err)

}

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
