package handlers

import (
	"log"
	"net/http"

	"github.com/Philip-21/bookings/internal/forms"

	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
	"golang.org/x/crypto/bcrypt"
)

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})

}

func (m *Repository) DisplaySignUp(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signup.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	//preventing session fixation attack by renewing the token
	_ = m.App.Session.RenewToken(r.Context())
	//_= m.App.Session.Cookie(r)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	firstname := r.Form.Get("firstname")
	lastname := r.Form.Get("lastname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

	form := forms.New(r.PostForm)
	form.Required("firstname", "lastname", "email", "password") //must be filled shows field cant be blank
	form.IsEmail("email")
	if !form.Valid() {
		//shows invalid email address based on the IsEmail format in forms.go
		//return an empty form and displays a message that
		render.Template(w, r, "signup.page.html", &models.TemplateData{
			Form: forms.New(nil),
		})
		return
	}

	user, err := m.DB.CreateUser(firstname, lastname, email, string(hashedPassword))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant fill sigup credentials!")
		http.Redirect(w, r, "/user/signup", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", user)
	m.App.Session.Put(r.Context(), "flash", "Signed up Successfully")
	//http redirect which directs to another page after the user fills a form,to prevent filling the form 2wice
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
		//Take user back to the main login page shoing an empty form to fill
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		//return back to the login form
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "User Doesnot exist")
		//m.App.Session.Put(r.Response.Context(), "error", "invalid cedentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
