package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var insertTimeout time.Duration = 3 * time.Second
var getTimeut time.Duration = 2 * time.Second
var updateTimeout time.Duration = 3 * time.Second
var deleteTimeout time.Duration = 2 * time.Second

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // "-" struct tag means, omit this from json serialization always
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	db *sql.DB
}

func NewMovieModel(db *sql.DB) MovieModel {
	return MovieModel{
		db: db,
	}
}

func (m MovieModel) Insert(ctx context.Context, movie *Movie) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	queryCtx, cancel := context.WithTimeout(ctx, insertTimeout)
	defer cancel()
	return m.db.QueryRowContext(queryCtx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(ctx context.Context, id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1
	`

	var movie Movie

	queryCtx, cancel := context.WithTimeout(ctx, getTimeut)
	defer cancel()

	err := m.db.QueryRowContext(queryCtx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(ctx context.Context, movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	queryCtx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()
	err := m.db.QueryRowContext(queryCtx, query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUpdateConflict
		default:
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM movies WHERE id = $1
	`
	queryCtx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()
	result, err := m.db.ExecContext(queryCtx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type MovieModelMock struct{}

func (m MovieModelMock) Insert(ctx context.Context, movie *Movie) error {
	return nil
}

func (m MovieModelMock) Get(ctx context.Context, id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieModelMock) Update(ctx context.Context, movie *Movie) error {
	return nil
}

func (m MovieModelMock) Delete(ctx context.Context, id int64) error {
	return nil
}
