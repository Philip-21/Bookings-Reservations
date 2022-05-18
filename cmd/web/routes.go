package main

import (
	"net/http"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//-------------------------------------defining a router---------------------------------------------------\\
func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer) // middleware.Recoverer absorb panics and prints the stack trace,panic occurs when a program cannot function
	router.Use(NoSurf)               //Nosurf adds CSRF protection to POST request
	router.Use(SessionLoad)

	//----------------------------------get request-----------------------------//
	router.Get("/user/login", handlers.Repo.ShowLogin)
	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)

	//rooms
	router.Get("/generals-quarters", handlers.Repo.Generals)
	router.Get("/majors-suite", handlers.Repo.Majors)

	//reservations and information
	router.Get("/search-availability", handlers.Repo.Availability)
	router.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	router.Get("/book-room", handlers.Repo.BookRoom)
	router.Get("/contact", handlers.Repo.Contact)
	router.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	router.Get("/make-reservation", handlers.Repo.Reservation)

	router.Get("/user/logout", handlers.Repo.Logout)

	//----------------------------post request-------------------------------------------------------//

	router.Post("/user/login", handlers.Repo.PostShowLogin)

	//reservations and information
	router.Post("/search-availability", handlers.Repo.PostAvailability)
	router.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	router.Post("/make-reservation", handlers.Repo.PostReservation)

	//gets the static files folder which contains the image
	fileServer := http.FileServer(http.Dir("./static/")) //.gets to the root of the application
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	///certain routes for the app are available to authenticated users(loged in users)
	router.Route("/admin", func(router chi.Router) {
		router.Use(Auth) //the login middleware
		router.Get("/dashboard", handlers.Repo.AdminDashboard)

		router.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		router.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		router.Get("/reservations-calender", handlers.Repo.AdminReservationsCalender)
		router.Post("/reservations-calender", handlers.Repo.AdminPostReservationsCalender)

		router.Get("/process-reservations/{src}/{id}", handlers.Repo.AdminProcessReservation)
		router.Get("/delete-reservations/{src}/{id}", handlers.Repo.AdminDeleteReservation)

		router.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)      //refers to the source and ID(/admin/reservations-new/src/id) of the new reservaton and all reservation page
		router.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation) //refers to the source and ID(/admin/reservations-new/src/id) of the new reservaton and all reservation page

	})
	return router
}
