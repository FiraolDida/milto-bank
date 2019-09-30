package Database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func connection() *sql.DB {
	db, err := sql.Open("postgres", "postgres://firaol:password@localhost/bankoffiraol?sslmode=disable");
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("You have successfully connect to the database")
	return db
}
