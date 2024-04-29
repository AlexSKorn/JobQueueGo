package routes

import (
	"testing"
)

func TestEnqueue(t *testing.T) {
	// Test to see if enqueue works in the best sense
}

func TestFailedEnqueue(t *testing.T) {
	// Try to enqueue but pass an error
}

func TestGetJob(t *testing.T) {
	// Enqueue a job and see if we can get it
}

func TestGetFailedJob(t *testing.T) {
	//don't enqueue a job or pass a bad id and make sure it errors properly
}

func TestDequeue(t *testing.T) {
	// enqueue and dequeue a job make sure it is the correct one
}

func TestFailedDequeue(t *testing.T) {
	// don't enqueue a job and try to dequeue
}

func TestConclude(t *testing.T) {
	// enqueue a job and conclude it
}

func TestFailedConclude(t *testing.T) {
	// enqueue, conclude a job, and then try to conclude it, it should not conclude
}
