package dbrepo

//---------------------initializing functions to be used by the  DatabaseRepo interface in repository.go-------------------\\
import (
	"context"
	"errors"
	"time"

	"github.com/Philip-21/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

//function available to the  DatabaseRepo interface,has the receiver postgresDBRepo

//connection for users
func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//initializing Reservation Repo so it can be used by the DatabaseRepointerface for swappping contents within the db and be used by the handlers
// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//using context to terminate a transaction e.g when a connection is lost fromthe user,or the user might close the page or browser
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	//speaking to the database, adding new entries to database
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
			end_date, room_id, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id` //returning id generates reservationID  to be used in the handlers

	err := m.DB.QueryRowContext(ctx, stmt, //running a query to return the reservation ID
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID) //scans the returning id into newID

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,	
			created_at, updated_at, restriction_id) 
			values
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select
			count(id)
		from
			room_restrictions
		where
			room_id = $1
			and $2 < end_date and $3 > start_date;`
	//$2 placehoder for start_date $3 placeholder for end_date

	//perform the query to the db
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range when searched by the user
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	//gets the id and room names from the rooms table, where the id in the rooms table is not in the room restrictions table to get a particular room available
	query := `
		select
			r.id, r.room_name
		from
			rooms r
		where r.id not in 
		(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
		`
	//querying the context
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			//scan room id and name
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1
`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

//
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query :=
		`select id , first_name, last_name, email, password, access_level, created_at, updated_at
	from users where id =$1`
	//using query role context cause we are getting a paticular row
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

//updates a user in the db
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5
	`
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

//authenticate a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//creating  variables  to hold the id and password of the authenticated user
	var id int
	var hashedPassword string

	//confirming the email from the database
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	//comparing and confirming the  password
	//matching the hashed password in the database to the testpassword inputed by user
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err

	}
	return id, hashedPassword, nil

}

//Return a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation

	query := `
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at,
	rm.id, rm.room_name
	from reservation r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close() //to avoid memory leaks
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			//scanning rows into room name and room id
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		//add the i variable to reservations
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation

	query :=
		`
	select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
 	rm.id, rm.room_name
	from reservation r 
	left join rooms rm on (r.room_id = rm.id)
	where processed = 0
	order by r.start_date asc
	` //where processed is a column that gives a default value of 0 when a customer makes a new reservation

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close() //to avoid memory leaks
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			//scanning rows into room name and room id
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		//add the i variable to reservations
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

//get a particular reservation
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query :=
		`
	   select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	   r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
	   rm.id, rm.room_name
	   from reservation r 
	   left join rooms rm on (r.room_id = rm.id)
	   where r.id= $1
	   `
	// left join rooms rm on (r.room_id = rm.id) refers to the foreign key
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

//making chnges and updates in the reservation Database
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
	`
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
