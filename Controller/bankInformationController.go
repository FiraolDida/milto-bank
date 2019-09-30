package Controller

import (
	"encoding/json"
	"fmt"
	"google.com/BankOfFiraol/Database"
	"google.com/BankOfFiraol/Model"
	"log"
	"net/http"
	"strconv"
	"time"
)

type data struct {
	UserDetail []Model.BankInformation
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("user not logged in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		email := cookie.Value
		fmt.Println(r.Method)

		if r.Method == http.MethodPost {
			m, _ := strconv.Atoi(r.FormValue("account"))
			accountNumber := int64(m)
			n, _ := strconv.Atoi(r.FormValue("amount"))
			amount := int64(n)

			depositInfo := Model.BankInformation{
				AccountNumber: accountNumber,
				Amount:        amount,
			}

			status, newAmount := Database.MakeDeposit(depositInfo)
			if status {
				fmt.Println(newAmount)
			} else {
				return
			}
		}
		bankInfo := []Model.BankInformation{}
		status, accountNumber := Database.GetUID(email, "deposit")
		if status {
			bankInfo = Database.GetDepositHistory(accountNumber)
			err := tpl.ExecuteTemplate(w, "deposit.html", bankInfo)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatalln(err)
			}
		}
	}
}

func Amount(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("user not logged in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		email := cookie.Value
		fmt.Println(r.Method)
		if r.Method == http.MethodPost {
			info := Model.BankInformation{}
			fmt.Println(info)

			err := json.NewDecoder(r.Body).Decode(&info)

			if err != nil {
				fmt.Println("error inside amount", err)
			}
			fmt.Println(info)

			status, amount := Database.GetAccountNumber(info.AccountNumber)

			if status {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				amountInfo := map[string]interface{}{
					"amount": amount,
				}
				err = json.NewEncoder(w).Encode(amountInfo)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
			}
		}
		status, amount := Database.GetUID(email, "amount")
		if status {
			err := tpl.ExecuteTemplate(w, "amount.html", amount)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatalln(err)
			}
		} else {
			fmt.Println("status is false")
		}
	}
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("user not logged in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		email := cookie.Value
		fmt.Println(r.Method)

		if r.Method == http.MethodPost {
			m, _ := strconv.Atoi(r.FormValue("account"))
			accountNumber := int64(m)
			n, _ := strconv.Atoi(r.FormValue("amount"))
			amount := int64(n)

			withdrawInfo := Model.BankInformation{
				AccountNumber: accountNumber,
				Amount:        amount,
			}

			status, currentAmount := Database.GetAccountNumber(withdrawInfo.AccountNumber)
			if status {
				if withdrawInfo.Amount < currentAmount {
					newAmount := currentAmount - withdrawInfo.Amount
					fmt.Println("new amount", newAmount, "current amount", currentAmount)
					t := time.Now()
					layout := "2 Jan 2006 15:04"
					date := t.Format(layout)
					flag := Database.UpdateAmount(newAmount, withdrawInfo.AccountNumber, date, "withdraw", withdrawInfo.Amount)

					if flag {
						fmt.Println("amount updated")
					} else {
						return
					}
				}
			}
		}
		bankInfo := []Model.BankInformation{}
		status, accountNumber := Database.GetUID(email, "withdraw")
		fmt.Println("an", status)
		if status {
			bankInfo = Database.GetWithdrawHistory(accountNumber)
			err := tpl.ExecuteTemplate(w, "withdraw.html", bankInfo)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatalln(err)
			}
		}
	}
}
