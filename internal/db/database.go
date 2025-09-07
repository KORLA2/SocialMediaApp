package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	duration, _ := time.ParseDuration(maxIdleTime)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}
