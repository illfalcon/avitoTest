package main

import (
	"log"

	"github.com/illfalcon/avitoTest/internal/interactions"
	"github.com/illfalcon/avitoTest/internal/sqlite"
)

func main() {
	db, cl, err := sqlite.StartService()
	if err != nil {
		log.Fatal(err, "could not connect to db")
	}
	is := interactions.NewMyInteractor(db)
	h := Handlers{&is}
	h.Start(cl)
}
