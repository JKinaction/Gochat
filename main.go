package main

import (
	"GoChat/model"
	"GoChat/routes"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	model.Initsql()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := routes.IntializeRoutes()
	r.Run(":" + port)
}
