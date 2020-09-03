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
	InfoColor    = "\033[1;34m%s\033[0m%s"
	NoticeColor  = "\033[1;36m%s\033[0m%s"
	WarningColor = "\033[1;33m%s\033[0m%s"
	ErrorColor   = "\033[1;31m%s\033[0m%s"
	SuccessColor = "\033[1;32m%s\033[0m%s"
	DebugColor   = "\033[0;36m%s\033[0m%s"
)

type TaskOptions struct {
	Text    string
	Execute func(wg *sync.WaitGroup, results chan TaskResult)
}

type TaskResult struct {
	Message string
	Error   bool
}

func Text(text string) TaskOption {
	return func(args *TaskOptions) {
		args.Text = text
	}
}

func Execute(executeFunc func(wg *sync.WaitGroup, results chan TaskResult)) TaskOption {
	return func(args *TaskOptions) {
		args.Execute = executeFunc
	}
}

type TaskOption func(*TaskOptions)

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
		fmt.Fprintln(os.Stdout, fmt.Sprintf("⌛ %s", t.Text))
	} else {
		s = spinner.New(spinner.CharSets[11], 150*time.Millisecond)
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
		fmt.Fprintln(os.Stdout, fmt.Sprintf(ErrorColor, "✗ ", result.Message))
	} else {
		fmt.Fprintln(os.Stdout, fmt.Sprintf(SuccessColor, "✓ ", result.Message))
	}
}
