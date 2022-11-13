package driver

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func UserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id int primary key , 
		first_name character varying(255),
		last_name character varying(255) ,
		email character varying(255) , 
		password character varying(60),
		created_at date,
		updated_at date,
	      access_level int default 1
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
		id int primary key,
		room_name character varying(225)
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
		id int primary key NOT NULL, 
		first_name character varying(255) NOT NULL, 
		last_name character varying(255) NOT NULL, 
		email character varying(225) NOT NULL,    
		phone character varying(255) NOT NULL, 
		start_date date, 
		end_date date, 
		room_id int ,
		processed int default 0
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
		id int primary key,
		start_date date, 
		end_date date, 
		room_id int,
		reservation_id int
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
