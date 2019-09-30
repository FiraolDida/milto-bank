package Controller

import (
	"fmt"
	"google.com/BankOfFiraol/Database"
	"google.com/BankOfFiraol/Model"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/pages/*"))
}

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		password := r.FormValue("password")

		newUser := Model.User{
			UID:       0,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		}

		status, _ := Database.CreateUser(newUser)
		if status {
			fmt.Println("You have successfully logged in")
			cookie, err := r.Cookie("session")
			if err != nil {
				cookie = &http.Cookie{
					Name:  "session",
					Value: newUser.Email,
				}
				http.SetCookie(w, cookie)
			}
			http.Redirect(w, r, "/amount", http.StatusSeeOther)
		} else {
			fmt.Println("Error creating user. Please try again!")
		}
	}
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		newUser := Model.User{
			Email:    email,
			Password: password,
		}

		fmt.Println(newUser)

		if Database.LoginUser(newUser) {
			fmt.Println("You have successfully logged in")
			cookie, err := r.Cookie("session")
			if err != nil {
				cookie = &http.Cookie{
					Name:  "session",
					Value: newUser.Email,
				}
				http.SetCookie(w, cookie)
			}
			http.Redirect(w, r, "/amount", http.StatusSeeOther)
		} else {
			fmt.Println("Invalid credentials")
			err := tpl.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatalln(err)
			}
		}

	} else {
		err := tpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalln(err)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("already logged out")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		cookie = &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}
