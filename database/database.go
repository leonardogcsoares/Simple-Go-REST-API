// Package database provides a wrapper to interact with the MySQL database instance.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Driver for mysql
)

const (
	dbConnection = "sql5116140:2sEyh7Q2RI@tcp(sql5.freemysqlhosting.net:3306)/sql5116140"
	tableName    = "tokenIds"
	tokenColumn  = "tokens"
	idColumn     = "ids"
)

// SaveIDToDb opens a connection to the database and stores the token-id string pair
func SaveIDToDb(tokenKey, id string) {
	db, err := sql.Open("mysql", dbConnection)
	checkErr(err)

	fmt.Println("Token key: "+tokenKey, "\tId: "+id)
	// Prepare the Insert statement used to insert row
	statement, err := db.Prepare("INSERT " + tableName + " SET " + tokenColumn + "=?," + idColumn + "=?")
	checkErr(err)

	// Execute the insertion, and return a response.
	res, err := statement.Exec(tokenKey, id)
	checkErr(err)
	fmt.Println(res)

	// CLose the connection
	db.Close()
}

// GetIDFromDb opens a connection to the database and given a tokenKey retrieves
// the id.
// Returns either the IMEI code or an "No IMEI found" message
func GetIDFromDb(tokenKey string) string {
	db, err := sql.Open("mysql", dbConnection)
	checkErr(err)

	rows, err := db.Query("SELECT * FROM " + tableName)
	checkErr(err)

	for rows.Next() {
		var tokenString string
		var imei string

		err = rows.Scan(&tokenString, &imei)
		fmt.Println(tokenString)
		fmt.Println(tokenKey)

		if tokenString == tokenKey {
			db.Close()
			return imei
		}

	}

	db.Close()
	return "No IMEI for the given Token"
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
