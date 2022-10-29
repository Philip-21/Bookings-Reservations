package driver

//this connects our application to a database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Philip-21/bookings/internal/config"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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
	db, err := sql.Open("pgx", dsn)
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
	log.Println("connected to postgres successfully")
	// driver, err := postgres.WithInstance(db, &postgres.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	// m, err := migrate.NewWithDatabaseInstance(
	// 	"./migrations",
	// 	"postgres", driver)
	// m.Up() //applies all up migrations
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	// log.Println("Migrations Successful")

	return dbConn, nil
}
