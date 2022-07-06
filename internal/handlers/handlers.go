package handlers

import (
	"net/http"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/driver"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
	"github.com/Philip-21/bookings/internal/repository"
	"github.com/Philip-21/bookings/internal/repository/dbrepo"
)

//always start a go function with a block letter so that it can be easily imported into anther directory e.g 	renders.Template(w, "home.page.html")
//handlers create response and receives request for the clients

//Repository helps to swap contents of our application with a minimal changes requiredin the code base
//Repo is the repository used by for new handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig       // repository for swapping contents within the app
	DB  repository.DatabaseRepo //repository for swapping contents within the db
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository { //driver.DB is the connection pool
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}
