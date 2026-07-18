package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kirki58/greenlight/m/internal/data"
)

var currentYear int = time.Now().Year()

// createMovieHandler for the "POST /v1/movies" endpoint.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieDto := data.MovieDto{}

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
	movie := movieDto.MapTo(nil)

	err = app.models.MovieRepository.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.createdResponse(w, r, movie, fmt.Sprintf("/v1/movies/%d", movie.ID))
}

// showMovieHandler for the "GET /v1/movies/:id" endpoint
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "Invalid ID in request path")
		return
	}

	mov, err := app.models.MovieRepository.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSONResponse(w, envelope{"movie": mov}, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateMovieHandler for the "PUT /v1/movies/:id" endpoint
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "Invalid ID in request path")
		return
	}

	mov, err := app.models.MovieRepository.Get(id)
	if  err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	movieDto := data.MovieDto{}
	app.readJSONRequest(w, r, &movieDto)

	if validationErors := app.uniValidator.ValidateBody(movieDto); validationErors != nil {
		app.failedValidationResponse(w, r, validationErors)
		return
	}

	_ = movieDto.MapTo(mov)
	if err := app.models.MovieRepository.Update(mov); err != nil{
		app.serverErrorResponse(w, r, err)
	}
	
	if err := app.writeJSONResponse(w, envelope{"movie": mov}, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// partialUpdateMovieHandler for the "PATCH /v1/movies/:id" endpoint
func (app *application) partialUpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "Invalid ID in request path")
		return
	}

	mov, err := app.models.MovieRepository.Get(id)
	if  err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	partialMovieDto := data.PartialMovieDto{}
	app.readJSONRequest(w, r, &partialMovieDto)

	if validationErors := app.uniValidator.ValidateBody(partialMovieDto); validationErors != nil {
		app.failedValidationResponse(w, r, validationErors)
		return
	}

	_ = partialMovieDto.MapTo(mov)
	if err := app.models.MovieRepository.Update(mov); err != nil{
		app.serverErrorResponse(w, r, err)
	}
	
	if err := app.writeJSONResponse(w, envelope{"movie": mov}, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteMovieHandler for the "DELETE /v1/movies/:id" endpoint
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "Invalid ID in request path")
		return
	}
	
	err = app.models.MovieRepository.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	if err := app.writeJSONResponse(w, nil, http.StatusNoContent, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
