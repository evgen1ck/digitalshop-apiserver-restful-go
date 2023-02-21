package postgres

//
//import (
//	"github.com/jmoiron/sqlx"
//	_ "github.com/lib/pq"
//	"time"
//)
//
//const defaultTimeout = 3 * time.Second
//
//type Postgres struct {
//	*sqlx.DB
//}
//
//func New(dsn string) (*Postgres, error) {
//	db, err := sqlx.Connect("postgres", "postgres://"+dsn+"?sslmode=disable")
//	if err != nil {
//		return nil, err
//	}
//
//	db.SetMaxOpenConns(25)
//	db.SetMaxIdleConns(25)
//	db.SetConnMaxIdleTime(5 * time.Minute)
//	db.SetConnMaxLifetime(2 * time.Hour)
//
//	return &Postgres{db}, nil
//}
