package main

import (
	"net/http"

	"github.com/Philip-21/bookings/internal/helpers"
	"github.com/justinas/nosurf"
)

// a middleware that writes to the console when somebody hits a page
// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{ //cookie generated to identify user when they visit a page or website
		HttpOnly: true,
		Path:     "/",              //cookie path which applies to the entire site
		Secure:   app.InProduction, //the app.InProduction refers to the variable defined in the main package,     production use true but in development set it to false
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// adding a middleware that tells the webserver to remember a state using sessions
// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// makes routes secure by only allowing logged in users to have access to certain parts,pages(routes) of the application
func Auth(next http.Handler) http.Handler {
	//calling the helpers func that requires a pointer to http.request as a parameter
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		helpers.ValidateToken(&http.Request{})

		if !helpers.IsAuthenticated(r) {
			//not authenticted
			session.Put(r.Context(), "error", "log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)

	})

}
