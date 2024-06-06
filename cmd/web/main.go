package main

import (
	"TODO/pkg/models/mysql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	todos    *mysql.TodoModel
	session *sessions.Session
	users *mysql.UserModel
}

func main() {
	addr := ":3000"
	dsn := "root:root@/todoApp?parseTime=true"

	infoLog, errorLog := initLogger()
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret data")
	flag.Parse()

	// To keep the main() function tidy put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	//creating instance app
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		todos:    &mysql.TodoModel{DB: db},
		session: session,
		users: &mysql.UserModel{DB: db},
	}
	//Initialize new http server struct
	serv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	//writing the message using new logger
	infoLog.Println("starting server on :3000", addr)
	er := serv.ListenAndServe()
	errorLog.Fatal(er)
}
