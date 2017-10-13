package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/josiahsleppy/golang-concurrent-server/collatzHelper"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../resources")))
	http.HandleFunc("/api/collatz", collatzHandler)
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
	value, elapsedTime := collatz.Collatz(numValue, concurrent == "true")
	fmt.Fprintf(w, "%d - First operation took %s to complete \n", value, elapsedTime)
}
