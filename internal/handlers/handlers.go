package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/driver"
	"github.com/Philip-21/bookings/internal/forms"
	"github.com/Philip-21/bookings/internal/helpers"
	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
	"github.com/Philip-21/bookings/internal/repository"
	"github.com/Philip-21/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi"
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

//---------------------------------------Reservation-----------------------------------\\

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		//redirects back to home page
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//getting a  room by the user
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	//storing them in a session
	m.App.Session.Put(r.Context(), "reservation", res)

	//convert time to string which can be used on the page and stored in the form
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	//putting it in a string map
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		//storing them in templates
		Form:      forms.New(nil), //this returns an empty form to fill when a user gets the reservation page
		Data:      data,           //an empty reservation the very first time the page is displayed
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //parsing form data
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//defining dates
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	//describing date format
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//coverting room id from a string to an int totally with the reservation table/model
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//making room being added to the templates
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	form := forms.New(r.PostForm) //r.PostForm is gotten from the url.Values in forms.go creates a form object and sends it back to the url

	//mandatory checks
	form.Required("first_name", "last_name", "email") //if any of this has an empty field then form will show an error
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			//createdform and data  will be saved in the template data
			Form: form,
			Data: data,
		})
		return
	}
	//parsing the reservation handlers to our Reservation db repo  which will make it speak to the database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return //return stops execution when theres an error
	}

	//building our room restriction model which will be linked to the reservation table
	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//throw the variable(reservation) into the session when we get to the reservation summary page
	//we pull the value out of the session send it to the template and display the information

	m.App.Session.Put(r.Context(), "reservation", reservation)

	//http redirect which directs to another page after the user fills a form,to prevent filling the form 2wice
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//get reservation [{m.App.Session.Put(r.Context(), "reservation", reservation)}] out of the session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//this methods prevents unauthorization or forging the url(/reservation-summary)only allowing the user to see his reservation as far as he's logged in
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//getting rid of the sessio,which removes data from the reservation
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{}) //Data format in templates-data.go
	data["reservation"] = reservation    //puting the reservation in the map

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	//render templates and parse data
	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//----------------------------------------Availability-----------------------------------------------------------//

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
} //w, r, "search-availability.page.html", &models.TemplateData{}   calls w http.ResponseWriter, r *http.Request, html string, td *models.TemplateData from the Rendermplate ta

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//these assigned values gets the forms that matches the inputs in the search avilability html page
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	//describing date format
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//connecting the database functions and gettin the dates
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//if the user searches for a room that is not available,
	//it stores he error in a session and redircts back to the page and prints no availability
	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{}) //Data format in templates-data.go
	data["rooms"] = rooms                //puting the rooms in the map

	//saves the dates and puts them in a session to be able to choose the rooms available
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)
	//render templates to choose a particuar room and parse data
	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

//json response for availability
type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	//call database function
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	//removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")
	//creating a header that tells the web browser that is receiving my response & what kind of response i'm sending it
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact renders the search availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

///-------------------Rooms------------------///

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// used to have next 6 lines
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	log.Println(err)
	//	m.App.Session.Put(r.Context(), "error", "missing url parameter")
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}

	//"explode" the URL into a slice of strings
	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	//we grab the third element of that slice (position 2, since slices start counting from 0), and parse that into an int.
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//Get the reservation variable where RoomID is located from a session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//updating the RoomID, and putting it back in the session
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	//creating a reservation in the book-room link which will take us to the rmake-reservation page
	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
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

//shows all reservations in Admin section
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	res, err := m.DB.AllReservations()
	if err != nil {
		log.Println(err)
		return
	}
	//creating slice for data to be used in the template
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "admin-all-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

//shows  reservation calender in Admin section
func (m *Repository) AdminReservationsCalender(w http.ResponseWriter, r *http.Request) {
	//assume there is no month or year specified
	now := time.Now()
	//specifying the month and year
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0) //date format for next month(y,m,d)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	//putting in a string map
	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	//current month and year Format
	stringMap["this_month"] = now.Format("02")
	stringMap["this_month_year"] = now.Format("2006")

	//getting the no of days and putting in a map
	//getting first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	//getting the rooms
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	//storing information about  the reservation and block from the calender in a data structure to be used in the template
	//range through the rooms variable
	for _, x := range rooms {
		//create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		//making sure there's one entry for every single day in the current month
		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			//reservation map for current day
			reservationMap[d.Format("2006-01-02")] = 0 //0 means room is available
			blockMap[d.Format("2006-01-02")] = 0
		}

		///get all restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomsByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		//if its a reservation we'll put the proper entry in the reservationMap likewise a block
		for _, y := range restrictions {
			if y.ReservationID > 0 {
				//its a reservation
				//loop through and put an entry for each of the dates
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					//has an entry to the date and reservationID which builds a link to the reservation
					reservationMap[d.Format("2006-01-02")] = y.ReservationID
				}
			} else {
				//its a block on the calender
				blockMap[y.StartDate.Format("2006-01-02")] = y.ID
			}
		}
		//gives a reservation or block map for every room
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		//storing the blockMap in the session,
		//this shows the blocks we are getting rid of and which ones are  new
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}

	render.Template(w, r, "admin-reservations-calender.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//makes reservation changes from the reservation calender(post request)
func (m *Repository) AdminPostReservationsCalender(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	//process blocks handles the logic to process blocks for things checked and unchecked

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	//loops through rooms
	for _, x := range rooms {
		//get the blockmap from the session(Get request) which contains the blocks for a given room at the point the calender was displayed  ,loop through the map
		//loop through the map
		//if we have an entry in the map that does not exist in our posted data,
		// and if restrictions id>0, then  its a block e need to remove

		//get the blockmap from the session(Get request) which contains the blocks for a given room at the point the calender was displayed  ,loop through the map
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		//loop through the map
		for name, value := range curMap {
			//ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				//only pay attention to values >0, and that are not in the form post
				//the rest are justb place holders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						//delete the retriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}

	}
	//handle new blocks
	for name, _ := range r.PostForm {
		//if the name of a posted element has a prefix add_block from the templates
		if strings.HasPrefix(name, "add_block") {
			//split the strings on the underscore
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			//getting the date
			t, _ := time.Parse("2006-01-02", exploded[3])
			//insert a new block
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%d&m=%d", year, month), http.StatusSeeOther)
}

//shows a reservation  in the each of the reservation section
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	//get the url and explode it to split it
	exploded := strings.Split(r.RequestURI, "/")
	//grab the id
	id, err := strconv.Atoi(exploded[4]) //[4] is the slice of the index of the url("/admin/reservations/new/ID")or("/admin/reservations/all/ID")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)
	//getting the source variable i.e /new
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap[src] = src

	//grabbing the year and month from query parameters to be used as hidden fields in the page
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	stringMap["month"] = month
	stringMap["year"] = year

	//getting reservation from the database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "admin-reservations-show.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil), //render the forms in a page
	})
}

//saves an edited reservation
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//get the url and explode it to split it
	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4]) //[4] is the slice of the index of the url("/admin/reservations/new/ID")or("/admin/reservations/all/ID")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap[src] = src

	//getting reservation from the database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	//update the database after changes are made
	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//check if theres a month and year value that are not empty strings
	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		//redirects to thw calender page with the correct month and year
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//marks  a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		log.Println(err)
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-#{src}", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//deletes a Reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	//http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-#{src}", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}
