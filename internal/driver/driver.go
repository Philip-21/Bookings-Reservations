package driver

//this connects our application to a database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Philip-21/bookings/internal/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

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
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	//ping troubleshoots connectivity
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)
	dbConn.SQL = db

	log.Println("Connected to postgres successfully")

	return dbConn, nil
}
