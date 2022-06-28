package handlers

import (
	"log"
	"net/http"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/driver"
	"github.com/Philip-21/bookings/internal/forms"
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

///--------------------Authentication-------------------////////

//shows login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil), //creating an empty form
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//preventing session fixation attack by renewing the token
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password") //must be filled shows field cant be blank
	form.IsEmail("email")              // only a valid email type
	if !form.Valid() {
		//Take user back to the main login page for an invalid  form
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		//return back to the login form
		m.App.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "user/login", http.StatusSeeOther)
	}
	///storing id in the session when authenticated  successfully
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in Successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther) //returns to home page after loggin successfully
}

//shows logout
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	//destroy the session
	_ = m.App.Session.Destroy(r.Context())
	//renew sesion token
	_ = m.App.Session.RenewToken(r.Context())
	//redirects to the login page
	http.Redirect(w, r, "user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})

}

//shows new reservations in Admin section
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	newres, err := m.DB.AllNewReservations()
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = newres
	render.Template(w, r, "admin-new-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}
