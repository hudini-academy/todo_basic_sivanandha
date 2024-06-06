package main

import (
	"net/http"
	"github.com/bmizerany/pat"
) 

func (app *application)routes() http.Handler {
	mux := pat.New()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	mux.Post("/tasks", app.session.Enable(http.HandlerFunc(app.addTask)))
	mux.Post("/tasks/delete",app.session.Enable(http.HandlerFunc(app.delTask)))
	mux.Get("/tasks/get", app.session.Enable(http.HandlerFunc(app.getTask)))
	mux.Post("/tasks/update", app.session.Enable(http.HandlerFunc(app.updateTask)))

	mux.Get("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", app.session.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", app.session.Enable(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}