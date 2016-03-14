package job

import "github.com/robxu9/kahinah/log"

// Job is a function that returns an error (if needed).
type Job func() error

var (
	// Queue is the job processing queue.
	Queue = make(chan Job, 250)
)

// ProcessQueue processes the next job.
func ProcessQueue() {
	next, ok := <-Queue
	if !ok { // channel closed
		return
	}

	err := next()
	if err != nil {
		log.Logger.Criticalf("job: unable to process job: %v", err)
	}
}
