package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/handlers"
	"github.com/zahnah/study-app/internal/helpers"
	"github.com/zahnah/study-app/internal/models"
	"github.com/zahnah/study-app/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to close database: %v\n", err)
		}
	}(db)
	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMain()

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}
	err = srv.ListenAndServe()
	log.Fatalln(err)
}

func run() (*sql.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("Can't create a template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	log.Println("connection to database..")
	os.LookupEnv("DATABASE_URL")
	fmt.Println("db: ", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
