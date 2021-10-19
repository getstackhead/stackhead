package routines

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
)

type Task struct {
	Name                string
	Run                 func(r RunningTask) error
	ErrorAsErrorMessage bool
}

type RunningTaskObj struct {
	Spinner *yacspin.Spinner
}

func (r *RunningTaskObj) PrintLn(text string) {
	if r.Spinner != nil {
		r.Spinner.Message(text)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, text+"\n")
	}
}
func (r *RunningTaskObj) SetSuccessMessage(text string) {
	if r.Spinner != nil {
		r.Spinner.StopMessage(text)
	}
}
func (r *RunningTaskObj) SetFailMessage(text string) {
	if r.Spinner != nil {
		r.Spinner.StopFailMessage(text)
	}
}

type RunningTask interface {
	PrintLn(text string)
	SetSuccessMessage(text string)
	SetFailMessage(text string)
}

func RunTask(task Task) {
	cfg := yacspin.Config{
		Frequency:         150 * time.Millisecond,
		CharSet:           yacspin.CharSets[11],
		Suffix:            " " + task.Name,
		SuffixAutoColon:   true,
		Message:           "",
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
	}

	var s *yacspin.Spinner
	if viper.GetBool("verbose") {
		_, _ = fmt.Fprintf(os.Stdout, "⌛ %s\n", task.Name)
	} else {
		spinner, err := yacspin.New(cfg)
		s = spinner
		if err == nil {
			s.Reverse()
			s.Start()
		}
	}
	runningTask := &RunningTaskObj{Spinner: s}
	err := task.Run(runningTask)
	if err != nil {
		if s != nil {
			if task.ErrorAsErrorMessage {
				s.StopFailMessage(err.Error())
			}
			s.StopFail()
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "✗ %s\n", err.Error())
		}
	} else {
		if s != nil {
			s.Stop()
		}
	}
}
