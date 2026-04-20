// Package main is the application entry point.
package main

import (
	"embed"
	"log"

	"gantt/internal/server"
)

//go:embed static/* templates/*.tmpl
var embeddedAssets embed.FS

func main() {
	s := server.New(embeddedAssets)
	log.Println("server started at http://localhost:8080")
	if err := s.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
