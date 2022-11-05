package driver

//this connects our application to a database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Philip-21/bookings/internal/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	_ "github.com/mattes/migrate/source/file"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// defining nature of connection pool
const maxOpenDbConn = 10 //max connections to db open at a given time
const maxIdleDbConn = 5  //connections that remain idle in the db pool
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres

func ConnectSQL(connect *config.Envconfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s  password=%s sslmode=%s",
		connect.Host, connect.Port, connect.DBName, connect.User, connect.Password, connect.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)
	dbConn.SQL = db

	//ping troubleshoots connectivity
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to postgres successfully")
	// drv, err := postgres.WithInstance(db, &postgres.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// m, err := migrate.NewWithDatabaseInstance("file://migrations/up", "postgres", drv)
	// if err != nil {
	// 	log.Fatal(err)
	// 	log.Println("error in migrations")
	// 	return nil, err
	// }
	// err = m.Up() //applies all up migrations
	// if err != nil {
	// 	if err == migrate.ErrNoChange || err == migrate.ErrLocked {
	// 		log.Printf("%v\n", err)

	// 	}

	// 	return nil, errors.Unwrap(err)

	// }
	// log.Println("Migrations Successful")

	return dbConn, nil
}
