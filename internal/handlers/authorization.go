package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Philip-21/bookings/internal/forms"
	"github.com/Philip-21/bookings/internal/helpers"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
	"golang.org/x/crypto/bcrypt"
)

///////----------------------Aut
type Sign struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Repository) DisplaySignUp(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signup.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {

	var cred Sign
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cred.Password), 8)
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	credentials := &models.Register{
		Email:    r.Form.Get(cred.Email),
		Password: r.Form.Get(string(hashedPassword)),
	}
	Regform := forms.New(r.PostForm)
	Regform.Required(cred.Email, string(hashedPassword))
	Regform.IsEmail(cred.Email)
	if !Regform.Valid() {
		resp := jsonResponse{
			Message: "invalid credentials",
		}
		//applies Indent to format the output
		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
		//render.Template(w, r, "home.page.html", &models.TemplateData{})

	}
	err = m.DB.CreateUser(*credentials)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant fill sigup credentials!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

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
