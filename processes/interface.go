package processes

import (
	"errors"

	"github.com/robxu9/kahinah/models"
)

type ProcessStatus int

const (
	ProcessOK ProcessStatus = iota
	ProcessFail
	ProcessRunning
	ProcessAborted
)

var (
	ErrEnded = errors.New("process: already ended")
	Mapping  = map[string]NewProcess{}
)

// NewProcess provides a method to intialise a new process with stored model
// data. This is because we can't directly store data, so we are forced to do
// this instead.
// State should be restored, including if the process has already started,
// ended, etc, and configuration is read from the platform passed in.
type NewProcess func(*models.ListStageProcess, map[string]interface{}) Process

// Process is a task in a stage, where its output helps to contribute to the
// result of the whole stage. It must be able to handle calls concurrently.
type Process interface {
	// Start starts the process. If it can't start, it returns an error.
	Start() error
	// Stop forces the process to stop. If it has already ended, it returns ErrEnded.
	Stop() error
	// Reset asks the process to clear all its existing state.
	// If the process is running, it forces a stop, erroring if it can't.
	Reset() error
	// Finalise writes out the final state and commits.
	// After this, changes are irreversible and the process has completed.
	// This method is called when the stage is about to end and everything has passed.
	// It is never called if the stage has failed at some point.
	Finalise() error

	// Status returns the current status of the process. Cronjobs may call this method
	// periodically to assess the current state of the process.
	Status() ProcessStatus

	// APIRequest takes in a user, action, and value, and returns a response
	// that will be JSONified and sent to the client (or an error that will be
	// returned instead). The user may be the system user (for background/daemon)
	// events.
	APIRequest(user *models.User, action string, value string) (interface{}, error)
	// APIMetadata returns general metadata about the current process state to
	// the client. (The list only cares whether it's OK or not, so the rest of
	// the data can go to the client.)
	APIMetadata(user *models.User) interface{}
}

// BuildProcess takes a ListStageProcess, reads configuration from the specified
// source defined in the parent list, then calls the corresponding NewProcess.
// Returns an error if the process is not available.
func BuildProcess(l *models.ListStageProcess) (Process, error) {
	// TODO: implement
	return nil, nil
}
