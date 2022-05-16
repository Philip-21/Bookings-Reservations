package dbrepo

import (
	"database/sql"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/Philip-21/bookings/internal/repository"
)

//main db repository for receiving contents
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB //db connection pool
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingsRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
