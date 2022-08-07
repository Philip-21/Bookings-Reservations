package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Philip-21/bookings/internal/config"
)

//this helpers file contains things that will be used in various parts of the application interms of error handling
var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	//write to the info log
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	//getting the trace of the error
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) //err.Error() is the error message debug.Stack()this gives the detailed information about the error that took place
	//the error log writes it to the terminal in development
	//in production the error log writes it in a log file ,put a directive for the user to check email ,see the error message in the log file and fix the error

	//writing error log to the terminal window
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

func IsSignup(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "email")
	return exist
}
