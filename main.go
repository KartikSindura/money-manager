package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KartikSindura/money/internal/app"
	"github.com/KartikSindura/money/internal/routes"
)

func main() {
	var port int = 8080
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("we are running on port %d", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
