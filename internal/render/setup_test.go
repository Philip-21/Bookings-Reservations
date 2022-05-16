package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

//http.ResponseWriter has interfaces of Header(),Writeheader(),Write()
//creating an interface that will satisfy http.ResponseWriter used in the render test
type myWriter struct{}

//creating a Header method
func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

//creating a writeheader method
func (tw *myWriter) WriteHeader(i int) {

}

//creating write method
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
