package config

import (
	"html/template"
	"log"

	"github.com/Philip-21/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

//configuration file will be accessed and imported into other parts of the application
//config file will only import standard libraries and not other files/dir of the application without any logic  to prevent import cycle and this makes the app not to compile

//holds the application configuration,this app config helps in adding items to it
//the main application repository for swapping contents
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template //new templates are stored
	InfoLog       *log.Logger                   //a standard library that writes information to log files
	ErrorLog      *log.Logger                   //this allows us to write logs and handle errors efficiently
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData //create a channel that will be avilable to all parts of the application
}
