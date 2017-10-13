package main

import (
	"time"

	"github.com/josiahsleppy/golang-concurrent-server/collatzHelper"
)

//Worker holds a channel to register itself as available on, a channel
//of work requests, and an "event" channel for shutdown
type worker struct {
	workerPool chan chan job
	jobChannel chan job
	quit       chan bool
}

//Response represents the results of the operation
type response struct {
	maxCount    int
	elapsedTime time.Duration
}

func newWorker(workerPool chan chan job) worker {
	return worker{workerPool, make(chan job), make(chan bool)}
}

//Start method starts the worker spinning, doing our actual
//computations as work requests come in
func (w worker) start() {
	go func() {
		for {
			//Register this worker onto the worker queue
			w.workerPool <- w.jobChannel

			//Select will block until one of the communications can proceed
			select {
			case job := <-w.jobChannel:
				//We have received a work request from the dispatcher
				maxCount, elapsedTime := collatz.Collatz(job.targetNumber, job.concurrent)
				response := response{maxCount, elapsedTime}
				//Send the results of our work on the results channel--our request
				//handler is listening for this
				job.results <- response
			case <-w.quit:
				//We have received a signal to stop
				return
			}
		}
	}()
}

//Stop signals the worker to stop listening for work requests
func (w worker) stop() {
	go func() {
		w.quit <- true
	}()
}
