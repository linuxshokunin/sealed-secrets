package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	compromised = "compromised"
)

// ScheduleJobWithTrigger creates a long-running loop that runs a job each
// loop
// returns a trigger function that runs the job early when called
func ScheduleJobWithTrigger(period time.Duration, job func()) func() {
	trigger := make(chan struct{})
	go func() {
		for {
			<-trigger
			job()
		}
	}()
	go func() {
		for {
			time.Sleep(period)
			trigger <- struct{}{}
		}
	}()
	return func() {
		trigger <- struct{}{}
	}
}

const (
	kubeChars     = "abcdefghijklmnopqrstuvwxyz0123456789-" // Acceptable characters in k8s resource name
	maxNameLength = 245                                     // Max resource name length is 253, leave some room for a suffix
)

func validateKeyPrefix(name string) (string, error) {
	if len(name) > maxNameLength {
		return "", fmt.Errorf("name is too long, must be shorter than %d, got %d", maxNameLength, len(name))
	}
	for _, char := range name {
		if !strings.ContainsRune(kubeChars, char) {
			return "", fmt.Errorf("name contains illegal character %c", char)
		}
	}
	return name, nil
}
