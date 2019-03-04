package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

func main() {
	connectstring := os.Getenv("CLOUDSQL_DB_USER") + ":" + os.Getenv("CLOUDSQL_DB_PASSWORD") + "@tcp(" + os.Getenv("CLOUDSQL_DB_HOST") + ")/cloudsqlclient"
	db, err := sql.Open("mysql", connectstring)
	if err != nil {
		fmt.Println(connectstring)
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	fmt.Println("pinging")
	if err != nil {
		log.Printf("ping hung")
		// if proxy not ready this is a good checkpoint
		log.Fatal(err.Error())
	}

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squareNum WHERE number = ?")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
					log.Printf(err.Error())
				}
			} else {
				log.Fatal(err.Error())
			}
		}
	}

	for {
		// Seed before random int
		rand.Seed(time.Now().UTC().UnixNano())

		time.Sleep(300 * time.Second)

		var squareNum int // we "scan" the result in here
		datNum := randInt(0, 24)

		// Query the square-number of 13
		err = stmtOut.QueryRow(datNum).Scan(&squareNum) // WHERE number = 13
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("The square number of %d is: %d\n", datNum, squareNum)
	}
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
