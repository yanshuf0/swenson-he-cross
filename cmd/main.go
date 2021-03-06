package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/yanshuf0/cross/pkg/data"
	"github.com/yanshuf0/cross/pkg/handle"
)

func main() {
	// Get path to assets from flag.
	assetsDir := flag.String("assetsDir", "../assets", "defines the path to asset directory")
	port := flag.String("port", ":8080", "sets the port to serve")
	flag.Parse()
	// Instantiate DB pool. I've done some abstractions here using interfaces
	// and custom structs (rather than using database/sql itself). Why?
	// This approach separates the db from the methods to allow for easily
	// "mockable" synchronous unit testing.
	db, err := data.NewDB()
	if err != nil {
		log.Fatalf("unable to connect to db, error: %v", err)
	}

	env := &handle.Env{DB: db, AssetsDir: assetsDir}
	// Initialize mux.
	m := handle.InitMux(env)

	// Allow Cors.
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(*port, handlers.CORS(originsOk, headersOk, methodsOk)(m)))
}
