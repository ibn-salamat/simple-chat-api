package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("server started")

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
