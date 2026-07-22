package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

const maxJSONRequestBodyBytes int64 = 1_048_576 // 1MB

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.

func (app *application) readIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Reads the JSON request from the given request body to the given destination type instance, Any issues during the decoding of the JSON request body will return an error with a clean public-facable message
func (app *application) readJSONRequest(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONRequestBodyBytes)

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	err := jsonDecoder.Decode(dst)

	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A dynamic error will be returned with the message "json: unknown field: <field-name>" when an unknown field is present in the json request body
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			return fmt.Errorf("The JSON request body contains an unknown field: %s", strings.TrimPrefix(err.Error(), "json: unknown field"))

		case errors.As(err, &maxBytesError):
			app.logger.Printf("%s issued a request whose body exceeds the maximum size limit, to %s endpoint", r.RemoteAddr, fmt.Sprintf("%s %s", r.Method, r.URL.String()))
			return errors.New("the request body exceeds the maximum size limit")

		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic,
		// rather than returning an error to our handler. At the end of this chapter
		// we'll talk about panicking versus returning errors, and discuss why it's an
		// appropriate thing to do in this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

			// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// We expect a single JSON object in the request for every endpoint
	if err := jsonDecoder.Decode(dst); !errors.Is(err, io.EOF) {
		return errors.New("Only a single JSON object is expected in the request body")
	}

	return nil
}

func (app *application) readQueriedRequest(w http.ResponseWriter, r *http.Request, dst any) map[string]string{
	err := app.schemaDecoder.Decode(dst, r.URL.Query())
	if err != nil {
		var multiErr schema.MultiError
		if errors.As(err, &multiErr) {
			var convErr schema.ConversionError
			var unknownKeyErr schema.UnknownKeyError
			var emptyFieldErr schema.EmptyFieldError

			var errs map[string]string
			for field, _ := range multiErr{
				switch{
				case errors.As(err, &convErr):
					errs[field] = fmt.Sprint("must be a valid %s value", convErr.Type.Kind())
				case errors.As(err, &unknownKeyErr):
					errs[field] = fmt.Sprintf("Is an unknown query parameter")
				case errors.As(err, &emptyFieldErr):
					errs[field] = fmt.Sprintf("Is an empty query parameter (not allowed)")
				}
			}
			return errs
			
		} else {
			// If it's not a MultiError, it's a developer bug (e.g. passed a non-pointer)
			panic(err)
		}
	}

	return nil
}