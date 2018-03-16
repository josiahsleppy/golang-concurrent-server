package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/josiahsleppy/golang-concurrent-server/collatz"
	// "github.com/josiahsleppy/helperFiles/collatz"
)

func main() {
	//Serve static files from a directory.
	http.Handle("/", http.FileServer(http.Dir("../resources")))
	//Or specify a custom handler for specific routes.
	http.HandleFunc("/api/collatz", collatzHandler)
	//ListenAndServe listens on a specified port and accepts new connections
	//in an infinite loop, spawning a new goroutine for each.
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func collatzHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	num := params["num"][0]
	concurrent := params["conc"][0]
	numValue, err := strconv.Atoi(num)
	if err != nil || numValue < 1 {
		io.WriteString(w, "Please enter a valid integer greater than zero.")
		return
	}
	// value, elapsedTime := collatz.Collatz(numValue)
	value, elapsedTime := collatz.Collatz(numValue, concurrent == "true")
	//http.ResponseWriter implements the io.Writer interface, which is why it can be
	//used here. Fprintf will call its Write method which writes to the response.
	fmt.Fprintf(w, "%d - Single operation took %s to complete \n", value, elapsedTime)
}
