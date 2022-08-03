package models

import "time"

// User is the user model
type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessLevel int    `json:"accesslevel"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int    `json:"id"`
	RoomName  string `json:"roomname"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int    `json:"id"`
	RestrictionName string `json:"Restrictionname"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	StartDate time.Time
	EndDate   time.Time
	RoomID    int `json:"roomID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Processed int  `json:"processed"`
	Room      Room //making a room field which includes all the info in Room struct which associated with the room id (not a compulsory field)
}

// RoomRestriction is the room restriction model
type RoomRestriction struct {
	ID            int `json:"id"`
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int `json:"roomid"`
	ReservationID int `json:"reservationid"`
	RestrictionID int `json:"restrictionid"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

///Mail Data holda an email message
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
