package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

//A global buffered channel that we can send work requests on. The dispatcher listens to it.
var jobQueue chan job

//Contains information relevant to the job to be run, and a channel to send back results.
type job struct {
	targetNumber int
	concurrent   bool
	results      chan response
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../resources")))
	http.HandleFunc("/api/collatz", collatzHandler)
	jobQueue = make(chan job, 100)
	//The number of workers to start up for best performance
	//really depends on the specific machine/expected load.
	dispatcher := newDispatcher(1000)
	dispatcher.run()
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

	results := make(chan response)
	//Create a new work item from the request values and send it to the job queue.
	work := job{numValue, concurrent == "true", results}
	fmt.Println("Received work request from client, sending to dispatcher")
	jobQueue <- work
	//Now we wait for the worker to send us the results.
	response := <-results
	fmt.Println("Received results from worker, sending response to client")
	fmt.Println()
	fmt.Fprintf(w, "%d - Single operation took %s to complete \n", response.maxCount, response.elapsedTime)
}
