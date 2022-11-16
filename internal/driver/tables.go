package driver

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Philip-21/bookings/internal/models"
)

func UserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id Serial primary key, 
		first_name character varying(255)  NOT NULL ,
		last_name character varying(255) NOT NULL,
		email character varying(255) NOT NULL, 
		password character varying(60) NOT NULL,
		access_level int default 1 NOT NULL,
		created_at timestamp(6) NOT NULL,
		updated_at timestamp(6) NOT NULL
		)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating users table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	log.Println("Users Created")
	return nil
}
func RoomTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS rooms(
		id serial primary key ,
		room_name character varying(225)  NOT NULL,
		created_at timestamp(6)  NOT NULL,
		updated_at timestamp(6)  NOT NULL
		
	)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating Reservation table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	log.Println("Rooms Created")
	return nil
}

func ReservationTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS reservation(
		id Serial primary key , 
		first_name character varying(255) NOT NULL, 
		last_name character varying(255) NOT NULL, 
		email character varying(225) NOT NULL,    
		phone character varying(255) NOT NULL, 
		start_date date  NOT NULL, 
		end_date date  NOT NULL, 
		room_id int  NOT NULL,
		created_at timestamp(6)  NOT NULL,
		updated_at timestamp(6)  NOT NULL,
		processed int default 0  NOT NULL
		)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating Reservation table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	log.Println("Reservation Created ")
	return nil
}

func RoomRestrictionTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS room_restrictions(
		id Serial primary key ,
		start_date date  NOT NULL, 
		end_date date  NOT NULL, 
		room_id int  NOT NULL,
		reservation_id int,
		created_at timestamp(6)  NOT NULL,
		updated_at timestamp(6)  NOT NULL
	)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating Reservation table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	log.Println("Room Restriction Created")
	return nil

}

func AlterTable(db *sql.DB) error {
	query := `ALTER TABLE reservation ADD FOREIGN KEY (room_id) REFERENCES rooms(id);
	          ALTER TABLE room_restrictions ADD FOREIGN KEY(room_id) REFERENCES rooms(id);
		    ALTER TABLE room_restrictions ADD FOREIGN KEY(reservation_id) REFERENCES reservation(id)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating Reservation table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	log.Println("Foreign keys created")
	return nil

}

// Insert into rooms(id, room_name) values(2, General Rooms)
func SeedDB(db *sql.DB, r *models.Room) error {

	query := `INSERT INTO rooms (id, room_name, created_at, updated_at) VALUES($1, $2, $3, $4);`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, r.ID, r.RoomName, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		log.Printf("Error %s when inserting row into rooms table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d rooms created ", rows)
	return nil
}
