package main

import (
	"fmt"
	"log"

	"github.com/illfalcon/avitoTest/internal/sqlite"
)

func main() {
	s, _ := sqlite.StartService()
	id, err := s.CreateChat("lkjbjhbjbkbbkjho", []int{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
