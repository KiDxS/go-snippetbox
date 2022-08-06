package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KiDxS/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         *mysql.UserModel
}

func main() {

	addrParameter := flag.String("addr", ":4000", "HTTP network address")
	databaseConfigParameter := flag.String("dsn", "no default value", "Sets the configuration of the database connection")
	secret := flag.String("secret", "9fa3530c21c13484bec7febb07666992", "Sets the secret to be used for the session")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*databaseConfigParameter)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, _ := newTemplateCache("./ui/html/")

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.SameSite = http.SameSiteStrictMode
	// Modifies the value of the application struct
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &mysql.SnippetModel{DB: db},
		session:       session,
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:         *addrParameter,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addrParameter)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

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
