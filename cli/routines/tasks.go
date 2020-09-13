package routines

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
)

const (
	// ErrorColor colors the first text red and the second text default color
	ErrorColor = "\033[1;31m%s\033[0m%s\n"
	// SuccessColor colors the first text green and the second text default color
	SuccessColor = "\033[1;32m%s\033[0m%s\n"
)

// TaskOptions is a configuration used to define the behaviour of tasks and processing functions
type TaskOptions struct {
	Text    string
	Execute func(wg *sync.WaitGroup, results chan TaskResult)
}

// TaskResult is the result of a task execution and expected to be returned into the respective channel
type TaskResult struct {
	// internal name
	Name    string
	Message string
	Error   bool
}

// Text assigns the given text to TaskOption.Text
func Text(text string) TaskOption {
	return func(args *TaskOptions) {
		args.Text = text
	}
}

// Execute assigns the given execution function to TaskOption.Execute
func Execute(executeFunc func(wg *sync.WaitGroup, results chan TaskResult)) TaskOption {
	return func(args *TaskOptions) {
		args.Execute = executeFunc
	}
}

// TaskOption is a single task setting
type TaskOption func(*TaskOptions)

// RunTask executes a task that can be configured using task options
func RunTask(options ...TaskOption) {
	t := &TaskOptions{
		Text: "Processing...",
	}
	for _, setter := range options {
		setter(t)
	}

	// Disable spinner when verbose mode is enabled as it does not like additional stdout messages
	var s *spinner.Spinner
	if viper.GetBool("verbose") {
		_, _ = fmt.Fprintf(os.Stdout, "⌛ %s\n", t.Text)
	} else {
		s = spinner.New(spinner.CharSets[11], 150*time.Millisecond)
		// s.ShowTimeElapsed = true
		s.Reverse()
		s.Suffix = " " + t.Text
		s.Start()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	results := make(chan TaskResult, 1)
	go t.Execute(&wg, results)
	wg.Wait()

	if !viper.GetBool("verbose") {
		s.Stop()
	}
	result := <-results
	if len(result.Message) == 0 {
		result.Message = t.Text
	}
	if result.Error {
		fmt.Fprintf(os.Stdout, ErrorColor, "✗ ", result.Message)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, SuccessColor, "✓ ", result.Message)
	}
}
