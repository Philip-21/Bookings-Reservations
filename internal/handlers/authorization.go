package handlers

import (
	"encoding/json"
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
	form.Required("email", "password") //must be filled shows field cant be blank
	form.IsEmail("email")
	if !form.Valid() {
		resp := jsonAuthorization{
			Message: "invalid credentials",
		}
		out, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	err = m.DB.CreateUser(firstname, lastname, email, string(hashedPassword))

	if err != nil {

		m.App.Session.Put(r.Context(), "error", "cant fill sigup credentials!")
		http.Redirect(w, r, "user/signup", http.StatusTemporaryRedirect)
	}
	m.App.Session.Put(r.Context(), "firstname", firstname)
	m.App.Session.Put(r.Context(), "lastname", lastname)
	m.App.Session.Put(r.Context(), "email", email)
	m.App.Session.Put(r.Context(), "password", password)
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
		//Take user back to the main login page for an invalid  form
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		resp := jsonAuthorization{
			Message: "user doesnot exist",
		}
		out, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
