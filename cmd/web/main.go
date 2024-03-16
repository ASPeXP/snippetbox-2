package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/aspexp/snippetbox-2/internal/models"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
}

func main(){

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil ))

	if err := godotenv.Load(); err != nil {
		logger.Error("load .env error")
	}
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DSN"), "MySQL Data Source Name")
	flag.Parse()

	
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// 	AddSource: true,
	// } ))
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil ))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
	}
	// mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView )
	// mux.HandleFunc("/snippet/create", app.snippetCreate)
	
	logger.Info("starting server", "addr", *addr)
	
	
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn )
	if err != nil {
		return nil, err 
	}
	if err = db.Ping(); err != nil {
		return nil, err 
	}
	return db, nil
}