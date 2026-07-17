package main

import (
	"context"
	"database/sql"
	"time"
)

type databaseCfg struct {
	dsn         string
	driver      string
	connTimeout time.Duration // Seconds
	pool struct{
		maxOpenConns int
		maxIdleConns int
		connMaxIdleTime time.Duration
	}
}

func (app *application) connectDb() (*sql.DB, error) {
	db, err := sql.Open(app.config.db.driver, app.config.db.dsn)
	if err != nil {
		return nil, err
	}

	// Pool configuration
	db.SetMaxOpenConns(app.config.db.pool.maxIdleConns)
	db.SetMaxIdleConns(app.config.db.pool.maxIdleConns)
	db.SetConnMaxIdleTime(app.config.db.pool.connMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.config.db.connTimeout))
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
