package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Philip-21/bookings/internal/models"
	"github.com/Philip-21/bookings/internal/render"
)

//----------------------------------------Availability-----------------------------------------------------------//

// Contact renders the search availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// Search Available Room Buttton in admin Dashboar template
// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
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
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}
	log.Println("Dates Inputed")
	//connecting the database functions and gettin the dates
	log.Println("Searching DB for all Rooms and getting its Dates....and comparing with inputed dates ")
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		log.Println("Room Not Available, Date is fully Booked")
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}
	log.Println("Searched the Database for All rooms By the Dates inputed")
	//if the user searches for a room that is not available,
	//it stores he error in a session and redircts back to the page and prints no availability
	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/admin/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{}) //Data format in templates-data.go
	data["rooms"] = rooms                //puting the rooms in the map

	//saves the dates and puts them in a session to be able to choose the rooms available
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	log.Println("Inpted Dates Created")
	m.App.Session.Put(r.Context(), "reservation", res)
	//render templates to choose a particuar room and parse data,
	//in the template, a link is generated  /admin/choose/room and a Get requests calls the choose room handler
	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

// The Make a Reservation button in the Admin Dashboard
// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	log.Println("Searching for Available Rooms Before Booking")
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonAvailability{
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
	log.Println("Searching DB for room by id And Dates.....")
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		log.Println("Room Not available, Date Currently Booked")
		// got a database error, so return appropriate json
		resp := jsonAvailability{
			OK:      false,
			Message: "Error querying database",
		}
		//applies Indent to format the output
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	log.Println("Search for room by Id and Dates completed")
	resp := jsonAvailability{
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
