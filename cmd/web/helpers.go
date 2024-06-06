package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	// "runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n", err.Error())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	http.Error(w, trace, http.StatusInternalServerError)
}

// The clientError helper sends status code and corresponding desc to user
// send responses like 400 "Bad Request" when there's a problem with the request that the user sent
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// also implement a notFound helper
// it is convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func initLogger() (*log.Logger, *log.Logger) {
	//create variable f for info
	infoFile, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//create variable e for error.log
	errFile, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//created logger using log.New() for writing information message
	//destination,prefix message,flags for additional info(three parameters)
	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

// creating openDB functionn wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
