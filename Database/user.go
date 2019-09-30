package Database

import (
	"database/sql"
	"fmt"
	"google.com/BankOfFiraol/Model"
	"log"
	"math/rand"
	"time"
)

func randomNumberGenerator(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func CreateUser(users Model.User) (bool, int) {
	db := connection()
	defer db.Close()

	id := time.Now().UnixNano()

	user := Model.User{}
	user.UID = id
	user.FirstName = users.FirstName
	user.LastName = users.LastName
	user.Email = users.Email
	user.Password = users.Password

	fmt.Println(user)

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		fmt.Println("fields are empty")
		return false, 0
	}

	_, err := db.Exec("INSERT INTO users (id, firstname, lastname, email, password) VALUES ($1, $2, $3, $4, $5)",
		user.UID, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		fmt.Println("error", err)
		return false, 0
	}
	log.Println("User successfully created")
	account_number := createAccountNumber(user.UID)
	accNum := int64(account_number)
	t := time.Now()
	layout := "2 Jan 2006 15:04"
	date := t.Format(layout)
	flag := UpdateAmount(0, accNum, date, "withdraw", 0)
	if flag {
		return true, account_number
	}
	return false, 0
}

func LoginUser(users Model.User) bool {
	db := connection()
	defer db.Close()

	user := Model.User{}
	user.Email = users.Email
	user.Password = users.Password

	fmt.Println(user)

	row := db.QueryRow("SELECT * FROM users WHERE email = $1 AND password = $2", user.Email, user.Password)

	loggedUser := Model.User{}
	err := row.Scan(&loggedUser.UID, &loggedUser.FirstName, &loggedUser.LastName, &loggedUser.Email, &loggedUser.Password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("invalid credentials", err)
		return false
	case err != nil:
		fmt.Println("database is nil", err)
		return false
	}
	fmt.Println(loggedUser.FirstName, loggedUser.LastName)
	return true
}

func createAccountNumber(id int64) int {
	db := connection()
	defer db.Close()

	t := time.Now()
	layout := "2 Jan 2006 15:04"

	for {
		accountNumber := randomNumberGenerator(100000, 999999)
		if accountNumberTester(accountNumber) {
			_, err := db.Exec("INSERT INTO bankInformation (id, accountNumber, amount) VALUES ($1, $2, $3)",
				id, accountNumber, 25)
			if err != nil {
				fmt.Println("error", err)
				return 0
			}
			log.Println("Account Information created successfully")
			DepositHistory(int64(accountNumber), 0, t.Format(layout))
			return accountNumber
			break
		}
	}

	return 0
}

func DepositHistory(id int64, amount int, dateValue string){
	db := connection()
	defer db.Close()
	fmt.Println(id)
	_, err := db.Exec("INSERT INTO depositHistory (accountnumber, amount, date) VALUES ($1, $2, $3)",
		id, amount, dateValue)
	if err != nil {
		fmt.Println("error", err)
	}
}

func accountNumberTester(accountNumber int) bool {
	fmt.Println("account number tester")
	db := connection()
	defer db.Close()
	row := db.QueryRow("SELECT accountNumber FROM bankInformation WHERE accountNumber = $1", accountNumber)
	accountInfo := Model.BankInformation{}
	err := row.Scan(&accountInfo.AccountNumber)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Error", err)
		return true
	case err != nil:
		fmt.Println("database is nil", err)
		return true
	}
	log.Println("account number match:", accountInfo.AccountNumber)
	return false
}
