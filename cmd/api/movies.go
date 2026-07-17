package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kirki58/greenlight/m/internal/data"
	"github.com/kirki58/greenlight/m/internal/dto"
)

var currentYear int = time.Now().Year()

// createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieDto := dto.MovieDto{}

	err := app.readJSONRequest(w, r, &movieDto)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	if validationErors := app.uniValidator.ValidateBody(movieDto); validationErors != nil {
		app.failedValidationResponse(w, r, validationErors)
		return
	}

	// Map to an actual Movie type
	movie := movieDto.Map()

	err = app.models.MovieRepository.Insert(&movie)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return
	}
	app.createdResponse(w, r, movie, fmt.Sprintf("/v1/movies/%d", movie.ID))
}

// showMovieHandler for the "GET /v1/movies/:id
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "Invalid ID in request path")
		return
	}

	dummyMovie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Awesome Movie",
		Runtime:   100,
		Year:      0,
		Version:   1,
	}

	if err := app.writeJSONResponse(w, envelope{"movie": dummyMovie}, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
