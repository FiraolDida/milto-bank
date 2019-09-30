package API

import (
	"encoding/json"
	"fmt"
	"google.com/BankOfFiraol/Database"
	"google.com/BankOfFiraol/Model"
	"net/http"
	"time"
)

func Amount(w http.ResponseWriter, r *http.Request)  {
	info := Model.BankInformation{}
	fmt.Println(info)

	err := json.NewDecoder(r.Body).Decode(&info)

	if err != nil {
		fmt.Println("error inside api amount" ,err)
	}
	fmt.Println(info)

	status, amount := Database.GetAccountNumber(info.AccountNumber)

	if status {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		amountInfo := map[string]interface{}{
			"amount" : amount,
		}
		err = json.NewEncoder(w).Encode(amountInfo)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
	}
}

func Deposit(w http.ResponseWriter, r *http.Request)  {
	info := Model.BankInformation{}
	fmt.Println(info)

	err := json.NewDecoder(r.Body).Decode(&info)

	if err != nil {
		fmt.Println("error inside api deposit" ,err)
	}
	fmt.Println(info)

	status, newAmount := Database.MakeDeposit(info)
	if status {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		new_amount := map[string]interface{}{
			"amount" : newAmount,
		}
		err = json.NewEncoder(w).Encode(new_amount)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
	} else {
		return
	}
}

func Withdraw(w http.ResponseWriter, r *http.Request)  {
	info := Model.BankInformation{}
	fmt.Println(info)

	err := json.NewDecoder(r.Body).Decode(&info)

	if err != nil {
		fmt.Println("error inside api withdraw" ,err)
	}
	fmt.Println(info)

	status, currentAmount := Database.GetAccountNumber(info.AccountNumber)
	fmt.Println("withdraw amount", info.Amount, "current amount", currentAmount)

	if status {
		if info.Amount < currentAmount{
			newAmount := currentAmount - info.Amount
			fmt.Println("new amount", newAmount, "current amount", currentAmount)
			t := time.Now()
			layout := "2 Jan 2006 15:04"
			date := t.Format(layout)
			flag := Database.UpdateAmount(newAmount, info.AccountNumber, date, "withdraw", info.Amount)

			if flag {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				new_amount := map[string]interface{}{
					"amount" : newAmount,
				}
				err = json.NewEncoder(w).Encode(new_amount)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	newUser := Model.User{}
	fmt.Println(newUser)

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Println("error inside api signup", err)
	}

	status, accountNumber := Database.CreateUser(newUser)
	if status {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		account_number := map[string]interface{}{
			"accountnumber" : accountNumber,
		}
		err = json.NewEncoder(w).Encode(account_number)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	currentUser := Model.User{}
	fmt.Println(currentUser)

	err := json.NewDecoder(r.Body).Decode(&currentUser)
	if err != nil {
		fmt.Println("error inside api signup", err)
	}

	status := Database.LoginUser(currentUser)
	if status {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "successfully logged in")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
