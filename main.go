package main

import (
	"google.com/milto-bank/API"
	"google.com/milto-bank/Controller"
	"net/http"
)

func main()  {
	//non api
	http.Handle("/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/login", Controller.Login)
	http.HandleFunc("/signup", Controller.Signup)
	http.HandleFunc("/amount", Controller.Amount)
	http.HandleFunc("/deposit", Controller.Deposit)
	http.HandleFunc("/withdraw", Controller.Withdraw)
	http.HandleFunc("/about", Controller.About)
	http.HandleFunc("/logout", Controller.Logout)

	//withAPI
	http.HandleFunc("/api/amount", API.Amount)
	http.HandleFunc("/api/deposit", API.Deposit)
	http.HandleFunc("/api/withdraw", API.Withdraw)
	http.HandleFunc("/api/signup", API.Signup)
	http.HandleFunc("/api/login", API.Login)

	http.ListenAndServe(":8082", nil)
}



