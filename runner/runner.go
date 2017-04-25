// Package runner provides for execution of Scoot work and retrieval
// of the results of that work. The core of this is the Runner interface.
// This package includes a standard runner as well as test and simulation
// implementations.
package runner

//go:generate mockgen -package=mocks -destination=mocks/runner.go github.com/scootdev/scoot/runner Service
//Note: using reflection syntax due to Service's nested interfaces. Unlike our other mocks, this requires its own package?

import (
	"bytes"
	"fmt"
	"time"
)

const NoRunnersMsg = "No runners available."
const RunnerBusyMsg = "Runner is busy"
const LoggingErrMsg = "Error initializing logging."

// A command, execution environment, and timeout.
type Command struct {
	// Command line to run
	Argv []string

	// Key-value pairs for environment variables
	EnvVars map[string]string

	// Kill command after timeout. Zero value is ignored.
	// TODO(dbentley): should this apply to the command, or include wall-time (including time spent
	// in a queue, or checking out the snapshot)?
	Timeout time.Duration

	// Runner can optionally use this to run against a particular snapshot. Empty value is ignored.
	SnapshotID string

	// TODO(jschiller): get consensus on design and either implement or delete.
	// Runner can optionally use this to specify content if creating a new snapshot.
	// Keys: relative src file & dir paths in SnapshotId checkout. May contain '*' wildcard.
	// Values: relative dest path=dir/base in new snapshot (if src has a wildcard, then dest path is treated as a parent dir).
	//
	// Note: nil and empty maps are different!, nil means don't filter, empty means filter everything.
	// SnapshotPlan map[string]string

	JobID  string
	TaskID string

	// ClientID is used to trace a task through the scheduler and a worker
	// ClientID string
}

func (c Command) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Command\n\tSnapshotID:\t%s\n\tArgv:\t%q\n\tTimeout:\t%v\n\tJobID:\t\t%s\n\tTaskID:\t\t%s\n",
		c.SnapshotID,
		c.Argv,
		c.Timeout,
		c.JobID,
		c.TaskID)

	if len(c.EnvVars) > 0 {
		fmt.Fprintf(&b, "\tEnv:\n")
		for k, v := range c.EnvVars {
			fmt.Fprintf(&b, "\t\t%s: %s\n", k, v)
		}
	}

	return b.String()
}

// Service allows starting/abort'ing runs and checking on their status.
type Service interface {
	Controller
	StatusReader

	// We need StatusEraser because Erase() shouldn't be part of StatusReader,
	// but we have it for legacy reasons.
	StatusEraser
}
