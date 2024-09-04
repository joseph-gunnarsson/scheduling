package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/joseph-gunnarsson/scheduling/api/handlers"
	"github.com/joseph-gunnarsson/scheduling/api/middleware"
	"github.com/joseph-gunnarsson/scheduling/api/routers"
	"github.com/joseph-gunnarsson/scheduling/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	conn := db.GetDBConnection()
	handler := handlers.NewBaseHandler(conn)
	mm := middleware.NewMiddlewareManager(conn)
	mux := routers.Routers(handler, mm)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
