package main

import (
	"database/sql"
	"flag"
	"github.com/dheepika6/LetsGoWebProgram/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {

	/** Extracting command level arguments */

	addr := flag.String("addr", ":4000", "Server port number")
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MYSQL data source string")
	flag.Parse()

	/** defining application level values */
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	/** Initialize database */

	db, err := OpenDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{logger: logger, snippets: &models.SnippetModel{DB: db}}

	logger.Info("The database connected successfully")

	mux := app.routes()
	logger.Info("starting server", slog.Any("addr", *addr))
	// log.Print("starting server at port", *addr)

	err = http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
