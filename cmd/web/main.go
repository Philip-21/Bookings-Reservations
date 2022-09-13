package main

import (
	"encoding/gob"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/driver"
	"github.com/Philip-21/bookings/internal/handlers"
	"github.com/Philip-21/bookings/internal/helpers"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// the main application function that runs the application
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err) //will stop the applicatio
	}
	defer db.SQL.Close()      //after the application stops ,db closes
	defer close(app.MailChan) //clossing the channel in the run function
	fmt.Println("starting mail listener..... ... ")
	listenForMail()

	//testing
	//msg:=models.MailData{
	//To    :  "john@gmail.com"
	//From    : hotel@gmail.com
	//Subject : Reservation
	//Content : ""
	//}
	//app.MailChan <-msg

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

// a function for testing
func run() (*driver.DB, error) {
	// what i am  going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	//gob.Register(models.Restriction{})
	gob.Register(models.User{})
	gob.Register(map[string]int{})
	gob.Register(models.TemplateData{})

	config.LoadConfig() //load the viper configuration
	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(config.Conf)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying.....")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache() //new templates which are stored in the createtemplatecache are defined as tc
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	//create a channel that will be avilable to all parts of the application
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false
	app.UseCache = false
	//setting up a logger to write to the terminal,helps in writing the client and server errors
	//info log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//error log
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//making it part of the app
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour              //session lasts for 24hours
	session.Cookie.Persist = true                  //storing sessions in cookies,cookie persists after the browser or window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode //specifying the site where the cookie applies to
	session.Cookie.Secure = app.InProduction       // insists the cookie is being crypted and the connection is from https //production use true but in development set it to false

	//defined for handlers to have access to
	app.Session = session

	app.TemplateCache = tc //defining the app as tc which stores template cache

	//application configurations
	repo := handlers.NewRepo(&app, db) //new repository and database configuration
	handlers.NewHandlers(repo)         //setting the repository in Newhandlers
	render.NewRenderer(&app)
	helpers.NewHelpers(&app) //pointing to a *config.AppConfig in the render dir

	return db, nil
}
