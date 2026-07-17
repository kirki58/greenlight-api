package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kirki58/greenlight/m/internal/data"
	_ "github.com/lib/pq"
)

// Declare a string containing the application version number. Later in the book we'll
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

const defaultDbDriver = "postgres"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   databaseCfg
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config       config
	logger       *log.Logger
	uniValidator *UniversalValidator
	models       data.Models
}

func main() {
	var cfg config
	parseCliFlags(&cfg)

	// Environment-specific configuruations
	if cfg.env == "development" && cfg.db.dsn == "" {
		cfg.db.dsn = os.Getenv("GREENLIGHT_DB_DSN_DEV")
	} else if cfg.env == "production" && cfg.db.dsn == "" {
		cfg.db.dsn = os.Getenv("GREENLIGHT_DB_DSN")
	}

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	uv, err := NewUniversalValidator()
	if err != nil {
		logger.Fatalln(err.Error())
	}
	app := &application{
		config:       cfg,
		logger:       logger,
		uniValidator: uv,
	}
	app.uniValidator.UseJSONTagNames()
	app.RegisterCustomValidations()
	db, err := app.connectDb()
	if err != nil {
		logger.Fatal("Could not establish a connection pool, ", err)
	}
	defer db.Close()

	logger.Println("Database connection pool established")
	app.setModels(db)

	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests
	// to the healthcheckHandler method (which we will create in a moment).
	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// Start the HTTP server.
	logger.Printf("starting %s server on %s\n", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func parseCliFlags(cfg *config) {
	// Declare an instance of the config struct.

	// application configuration from command-line flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Db configuration
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "Database connection string to be used, It is recommended to pass this as an environment variable instead which is the default")
	flag.StringVar(&cfg.db.driver, "db-driver", defaultDbDriver, "Driver name to be used for DB connection")
	connTimeoutFlag := *flag.String("db-conn-timeout", "5s", "First ping response timeout for the Database")
	flag.IntVar(&cfg.db.pool.maxOpenConns, "db-max-open-conn", 25, "Maximum open database connections in the pool")
	flag.IntVar(&cfg.db.pool.maxIdleConns, "db-max-idle-conn", 15, "Maximum number of database connections sitting idle in the pool (can't exceed db-max-open-conn)")
	maxIdleTimeFlag := *flag.String("db-max-idle-time", "15m", "Maximum time duration of connections sitting idle in the pool before they are terminated")
	flag.Parse()

	parsedDuration, err := time.ParseDuration(connTimeoutFlag)
	if err != nil {
		log.Fatal("db-conn-timeout should be a parsable time expression: \"5s\", \"10m\", \"1d\" etc.")
	}
	cfg.db.connTimeout = parsedDuration

	parsedMaxIdle, err := time.ParseDuration(maxIdleTimeFlag)
	if err != nil {
		log.Fatal("db-max-idle should be a parsable time expression: \"5s\", \"10m\", \"1d\" etc.")
	}
	cfg.db.pool.connMaxIdleTime = parsedMaxIdle
}

func (app *application) setModels(db *sql.DB){
	app.models = data.Models{
		MovieRepository: data.NewMovieModel(db),
	}
}