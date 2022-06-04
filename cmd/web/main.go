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
	"github.com/spf13/viper"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//the main application function that runs the application
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

//a function for testing
func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})
	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s ", err)
	}

	dbHost, ok := viper.Get("DB_host").(string)
	dbPort, ok := viper.Get("DB_port").(string)
	dbUser, ok := viper.Get("DB_user").(string)
	dbPassword, ok := viper.Get("DB_password").(string)
	dbdatabase, ok := viper.Get("DB_name").(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	fmt.Printf("viper : %s = %s\n", "host", dbHost)
	fmt.Printf("viper : %s = %s\n", "port", dbPort)
	fmt.Printf("viper : %s = %s\n", "dbname", dbdatabase)
	fmt.Printf("viper : %s = %s\n", "user", dbUser)
	fmt.Printf("viper : %s = %s\n", "password", dbPassword)

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=%s   port=%s dbname=%s  user=%s  password=philippians")
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
	//making it part of the app
	app.InfoLog = infoLog

	//error log
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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
