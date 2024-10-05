package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fabiobap/go-bnb/internal/config"
	"github.com/fabiobap/go-bnb/internal/driver"
	"github.com/fabiobap/go-bnb/internal/handlers"
	"github.com/fabiobap/go-bnb/internal/helpers"
	"github.com/fabiobap/go-bnb/internal/models"
	"github.com/fabiobap/go-bnb/internal/render"
	"github.com/joho/godotenv"

	"github.com/alexedwards/scs/v2"
)

const HTTP_PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	fmt.Printf("Starting email listener")
	listenForEmail()

	fmt.Printf("Starting application on port: %s", HTTP_PORT)

	srv := &http.Server{
		Addr:    HTTP_PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	loadVariables()

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	log.Println("Connecting to database...")

	var dbCreds = models.DBData{}
	setDBData(&dbCreds)

	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbCreds.DBHost,
		dbCreds.DBPort,
		dbCreds.DBName,
		dbCreds.DBUser,
		dbCreds.DBPass,
		dbCreds.DBSSL)

	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

func loadVariables() {
	isProductionStr := os.Getenv("IS_PRODUCTION")
	if isProductionStr == "" {
		isProductionStr = "false"
	}

	isProduction, err := strconv.ParseBool(isProductionStr)
	if err != nil {
		log.Fatalf("Error converting IS_PRODUCTION key to bool: %v", err)
	}

	app.InProduction = isProduction

	useCacheStr := os.Getenv("USE_CACHE")
	if useCacheStr == "" {
		useCacheStr = "false"
	}

	useCache, err := strconv.ParseBool(useCacheStr)
	if err != nil {
		log.Fatalf("Error converting USE_CACHE key to bool: %v", err)
	}

	app.UseCache = useCache
}

func setDBData(db *models.DBData) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "bookings"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASSWORD")

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbSSL := os.Getenv("DB_SSL")
	if dbSSL == "" {
		dbSSL = "disable"
	}

	db.DBHost = dbHost
	db.DBName = dbName
	db.DBUser = dbUser
	db.DBPass = dbPass
	db.DBPort = dbPort
	db.DBSSL = dbSSL
}
