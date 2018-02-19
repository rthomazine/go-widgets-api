package main

import (
	"fmt"
	"log"
	"time"
	"database/sql"

	_ "github.com/lib/pq"
)

type userToken struct {
	id        int64
	usercode  string
	usertoken string
	created   time.Time
	expires   time.Time
}

type users struct {
	id       int64
	name     string
	gravatar string
	created  time.Time
	updated  time.Time
}

type widgets struct {
	id        int64
	name      string
	color     string
	price     string
	inventory int64
	created   time.Time
	updated   time.Time
}

func openDB() *sql.DB {
	connString := 	fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable",
		"widgets_db", "localhost", "widgetsuser", "ChangeMe123.")

	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Error connecting to the database: ", err)
	}
	log.Println("Connected to database")

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func createUserToken(usertoken userToken) error {
	sqlStatement := "INSERT INTO widgets.user_token (usercode, usertoken, expires) " +
		"VALUES ($1,$2,$3)"

	_, err := db.Exec(sqlStatement, usertoken.usercode, usertoken.usertoken, usertoken.expires)
	if err != nil {
		log.Printf("Failed to insert user token %s: ", usertoken.usercode, err)
		return err
	}

	return nil
}

func queryUserToken(token string) (*userToken, error) {
	usertoken := userToken{}

	sqlQuery := "SELECT id, usercode, usertoken, created, expires " +
		"FROM widgets.user_token " +
		"WHERE usertoken = $1"
	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Printf("Failed to query user token [%s]: ", usertoken.usertoken, err)
		return nil, err
	}

	//row := db.QueryRow(sqlQuery, token)
	//err := row.Scan(&usertoken.id,
	//	&usertoken.usercode,
	//	&usertoken.usertoken,
	//	&usertoken.created,
	//	&usertoken.expires)
	stmt.QueryRow(token).Scan(&usertoken)
	if err != nil {
		return nil, err
	}

	return &usertoken, nil
}
