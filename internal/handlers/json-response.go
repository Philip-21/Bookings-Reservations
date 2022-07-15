package handlers

//json response for availability
type jsonAvailability struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type jsonAuthorization struct {
	Message string `json:"message"`
}
