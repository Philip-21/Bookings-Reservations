package main

import (
	"net/http"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// -------------------------------------defining a router---------------------------------------------------\\
func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https//*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept", "Authorization", " Accept-Encoding",
			"Content-Type", "Connection", " Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header", "Token", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Recoverer) // middleware.Recoverer absorb panics and prints the stack trace,panic occurs when a program cannot function
	router.Use(NoSurf)               //Nosurf adds CSRF protection to POST request
	router.Use(SessionLoad)          // SessionLoad loads and saves the session on every request

	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)
	router.Get("/contact", handlers.Repo.Contact)

	//----------------------------------Authorization Requests-----------------------------//
	router.Get("/user/signup", handlers.Repo.DisplaySignUp)
	router.Get("/user/login", handlers.Repo.ShowLogin)
	router.Get("/user/logout", handlers.Repo.Logout)
	router.Post("/user/login", handlers.Repo.PostShowLogin)
	router.Post("/user/signup", handlers.Repo.SignUp)

	//gets the static files folder which contains the image
	fileServer := http.FileServer(http.Dir("./static/")) //.gets to the root of the application
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	///certain routes for the app are available to authenticated users(loged in users)
	router.Route("/admin", func(router chi.Router) {
		router.Use(Auth) //the login middleware
		router.Get("/dashboard", handlers.Repo.AdminDashboard)

		//rooms
		router.Get("/generals-quarters", handlers.Repo.Generals)
		router.Get("/majors-suite", handlers.Repo.Majors)
		router.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

		//router.Get("/reservations-calender", handlers.Repo.AdminReservationsCalender)
		//router.Post("/reservations-calender", handlers.Repo.AdminPostReservationsCalender)

		//------reservations and information----------------//
		router.Get("/search-availability", handlers.Repo.Availability)
		router.Post("/search-availability", handlers.Repo.PostAvailability)
		router.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
		router.Get("/book-room", handlers.Repo.BookRoom)
		router.Post("/make-reservation", handlers.Repo.PostReservation)
		router.Get("/make-reservation", handlers.Repo.Reservation)
		router.Get("/reservation-summary", handlers.Repo.ReservationSummary)
		router.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		router.Get("/reservations-all", handlers.Repo.AdminAllReservations)

		//handling reservations
		router.Get("/process-reservations/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		router.Get("/delete-reservations/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		router.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation) //refers to the source and ID(/admin/reservations-new/src/id) of the new reservaton and all reservation page
		router.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation) //refers to the source and ID(/admin/reservations-new/src/id) of the new reservaton and all reservation page

	})
	return router
}
