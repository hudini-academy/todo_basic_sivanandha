package main

import (
	//"fmt"
	// "errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	//"strings"
	//"strconv"
	"TODO/pkg/models"
)

// created handler for home
// changed with *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Latest() assigining to a variable s
	s, err := app.todos.Latest()
	if err != nil {
		app.serverError(w, err)
		log.Println(err)
		return
	}

	// storing all templates to a files
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	//create template.parsefile for reading the template
	ts, err := template.ParseFiles(files...)
	//if there is error we will log print the error to the user
	if err != nil {
		//change to errorLog that we created in main
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error1", 500)
		return
	}

	err = ts.Execute(w, struct {
		Tasks []*models.Todos
		Flash string
	}{
		Tasks: s,
		Flash: app.session.PopString(r, "flash"),
	})

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error2", 500)
	}

}

// here also changing the function
func (app *application) addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	// inserting the values from form
	taskName := r.FormValue("text")
	taskDesc := r.FormValue("message")
	fmt.Println(taskName[:8])
	if strings.Contains(taskName, "Special") {
		fmt.Println("Here")
		_, err := app.specials.Insert(taskName, taskDesc, "10")
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "internal server error", 500)
		}
	}else {
		if strings.TrimSpace(taskName) == "" && strings.TrimSpace(taskDesc) == "" {
			app.session.Put(r, "flash", "This field cannot be blank")
		} else if utf8.RuneCountInString(taskName) > 100 {
			app.session.Put(r, "flash", "This field is too long (maximum is 100 characters")
		} else {
			_, err := app.todos.Insert(taskName, taskDesc, "365")
			if err != nil {
				app.serverError(w, err)
				return
			} else {
				app.session.Put(r, "flash", "Task successfully created!")
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) specialTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	st, errr := app.specials.Latest()
	if errr != nil {
		app.serverError(w, errr)
		log.Println(errr)
		return
	}
	files := []string{
		"./ui/html/special.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	err = ts.Execute(w, struct {
		Sp []*models.Special
		Flash string
	}{
		Sp: st,
		Flash: app.session.PopString(r, "flash"),
	})
	if err != nil{
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
// create function for delete special task
func (app *application) delSpecialTask(w http.ResponseWriter, r *http.Request){
	delSp := r.URL.Query().Get("id")
	delId, err := strconv.Atoi(delSp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
	}
	// assigning delete to a variable 'er'
	er := app.specials.Delete(delId)
	if er != nil {
		log.Println(er.Error())
		http.Error(w, "internal server error", 500)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// function for get one task based on their id
func (app *application) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.todos.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "%v", s)
}

// function for delete and changed the function with method *application
func (app *application) delTask(w http.ResponseWriter, r *http.Request) {
	del := r.URL.Query().Get("id")
	delId, err := strconv.Atoi(del)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
	}
	// assigning delete to a variable 'er'
	er := app.todos.Delete(delId)
	if er != nil {
		log.Println(er.Error())
		http.Error(w, "internal server error", 500)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// created function for update task
func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	upd := r.URL.Query().Get("id")
	upId, err := strconv.Atoi(upd)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
	}
	updateValue := r.FormValue("update")
	if strings.TrimSpace(updateValue) == "" {
		app.session.Put(r, "flash", "This field cannot be blank")
	} else {
		_, er := app.todos.Update(upId, updateValue)
		if er != nil {
			app.serverError(w, er)
			return
		} else {
			app.session.Put(r, "flash", "Task successfully Updated!!!!")
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	ts.Execute(w, app.session.PopString(r, "flash"))

}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("name")
	userEmail := r.FormValue("email")
	userPass := r.FormValue("password")
	err := app.users.Insert(userName, userEmail, userPass)
	if strings.TrimSpace(userName) == "" && strings.TrimSpace(userEmail) == "" {
		app.session.Put(r, "flash", "This field cannot be blank")
	} else if err != nil {
		app.errorLog.Println(err.Error())
		app.session.Put(r, "flash", "User already exist")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "flash", "Signup successful!")
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	ts.Execute(w, app.session.PopString(r, "flash"))

}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	userPass := r.FormValue("password")
	isUser, err := app.users.Authenticate(userEmail, userPass)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	if isUser {
		app.session.Put(r, "Authentication", true)
		app.session.Put(r, "flash", "Login Successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.session.Put(r, "flash", "Login failed")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.session.Put(r, "Authentiaction", false)

	}
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "id")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/user/login", 303)
}
