package main

import (
	"log"
	"movies/core/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
