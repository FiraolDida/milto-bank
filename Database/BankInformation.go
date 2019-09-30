package Database

import (
	"database/sql"
	"fmt"
	"google.com/milto-bank/Model"
	"time"
)

func MakeDeposit(infos Model.BankInformation) (bool, int64) {
	db := connection()
	defer db.Close()

	t := time.Now()
	layout := "2 Jan 2006 15:04"

	info := Model.BankInformation{}
	info.AccountNumber = infos.AccountNumber
	info.Amount = infos.Amount
	info.Date = t.Format(layout)

	fmt.Println(info.UID, " ", info.AccountNumber, " ", info.Amount, " Date", info.Date)

	status, currentAmount := GetAccountNumber(info.AccountNumber)

	if status {
		newAmount := info.Amount + currentAmount
		if UpdateAmount(newAmount, info.AccountNumber, info.Date, "deposit", info.Amount) {
			fmt.Println("Amount updated")
			return true, newAmount
		}
	}

	return false, 0
}

func GetAccountNumber(account int64) (bool, int64) {
	db := connection()
	defer db.Close()
	row := db.QueryRow("SELECT * FROM bankinformation WHERE accountNumber = $1", account)
	userAccount := Model.BankInformation{}
	err := row.Scan(&userAccount.UID, &userAccount.AccountNumber, &userAccount.Amount)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Account number not found1", err)
		return false, 0
	case err != nil:
		fmt.Println("Account number not found2", err)
		return false, 0
	}
	fmt.Println(userAccount.UID, " ", userAccount.AccountNumber, " ", userAccount.Amount)
	return true, userAccount.Amount
}

func UpdateAmount(newAmount int64, AccountNumber int64, dateValue string, from string, currentAmount int64) bool {
	db := connection()
	defer db.Close()
	fmt.Println(dateValue)

	_, err := db.Exec("UPDATE bankinformation SET amount = $1 WHERE accountnumber = $2", newAmount, AccountNumber)
	if err != nil {
		fmt.Println("Error updating account", err)
		return false
	}
	fmt.Println("Cash successfully deposited")
	t := time.Now()
	layout := "2 Jan 2006 15:04"
	if from == "deposit"{
		DepositHistory(AccountNumber, int(currentAmount), t.Format(layout))
		return true
	}
	WithdrawHistory(AccountNumber, int(currentAmount), t.Format(layout))
	return true
}

func WithdrawHistory(accountnumber int64, amount int, dateValue string){
	db := connection()
	defer db.Close()
	fmt.Println(accountnumber)
	_, err := db.Exec("INSERT INTO withdrawhistory (accountnumber, amount, date) VALUES ($1, $2, $3)",
		accountnumber, amount, dateValue)
	if err != nil {
		fmt.Println("error", err)
	}
}

func GetUID(email string, from string) (bool, int64) {
	db := connection()
	defer db.Close()
	row := db.QueryRow("SELECT id FROM users WHERE email = $1", email)
	userInfo := Model.User{}
	err := row.Scan(&userInfo.UID)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Account number not found", err)
		return false, 0
	case err != nil:
		fmt.Println("Account number not found", err)
		return false, 0
	}
	fmt.Println(userInfo.UID)
	status, amount := GetAccountInfoViaUID(userInfo.UID, from)
	fmt.Println(status, amount)
	if status {
		return true, amount
	}
	return false,0
}

func GetAccountInfoViaUID(id int64, from string) (bool, int64) {
	db := connection()
	defer db.Close()
	row := db.QueryRow("SELECT * FROM bankinformation WHERE id = $1", id)
	userAccount := Model.BankInformation{}
	err := row.Scan(&userAccount.UID, &userAccount.AccountNumber, &userAccount.Amount)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Account number not found", err)
		return false, 0
	case err != nil:
		fmt.Println("Account number not found", err)
		return false, 0
	}
	fmt.Println(userAccount.UID, " ", userAccount.AccountNumber, " ", userAccount.Amount)
	if from == "amount" {
		return true, userAccount.Amount
	} else if from == "deposit"{
		return true, userAccount.AccountNumber
	} else if from == "withdraw" {
		return true, userAccount.AccountNumber
	}
	return false, 0
}

func GetDepositHistory(accountNumber int64) []Model.BankInformation {
	db := connection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM deposithistory WHERE accountnumber = $1", accountNumber)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	bankInfo := make([]Model.BankInformation, 0)

	for rows.Next() {
		info := Model.BankInformation{}
		err := rows.Scan(&info.UID, &info.AccountNumber, &info.Amount, &info.Date)
		if err != nil {
			panic(err)
		}
		bankInfo = append(bankInfo, info)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	return bankInfo
}

func GetWithdrawHistory(accountNumber int64) []Model.BankInformation {
	db := connection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM withdrawhistory WHERE accountnumber = $1", accountNumber)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	bankInfo := make([]Model.BankInformation, 0)

	for rows.Next() {
		info := Model.BankInformation{}
		err := rows.Scan(&info.UID, &info.AccountNumber, &info.Amount, &info.Date)
		if err != nil {
			panic(err)
		}
		bankInfo = append(bankInfo, info)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	return bankInfo
}

