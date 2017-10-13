package main

//Dispatcher manages workers and the job queue
type dispatcher struct {
	//A pool of worker's channels that are registered with the dispatcher
	workerPool chan chan job
	numWorkers int
}

func newDispatcher(numWorkers int) *dispatcher {
	//Our worker pool will be buffered to the same number as there will be workers
	workerPool := make(chan chan job, numWorkers)
	return &dispatcher{workerPool, numWorkers}
}

func (d *dispatcher) run() {
	//Allocate the specified allowed number of workers
	for i := 0; i < d.numWorkers; i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}
	//Start the dispatcher spinning in a goroutine so this method can return
	go d.dispatch()
}

//Dispatch, once started, will spin until the program exits, incessantly listening for work
//requests and sending them off to an available worker
func (d *dispatcher) dispatch() {
	for {
		//Wait for a work request from our main request handler
		job := <-jobQueue
		//Once received, send off a goroutine so we can start listening
		//for work almost immediately
		go func() {
			//Wait for a worker to be available--we'll know when it's job channel gets
			//registered to the dispatcher's worker pool. This will block until someone's
			//available, but since we're in a goroutine, the enclosing function doesn't block.
			jobChannel := <-d.workerPool
			//Dispatch the job to the worker job channel
			jobChannel <- job
		}()
	}
}
