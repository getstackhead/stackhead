package routines

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/spf13/viper"
)

type Task struct {
	Name                string
	Run                 func(r *Task) error
	ErrorAsErrorMessage bool
	Spinner             *ysmrr.Spinner
	TaskRunner          *TaskRunner

	// disabled tasks are skipped silently
	Disabled bool

	SubTasks                   []Task
	IgnoreSubtaskErrors        bool
	RunAllSubTasksDespiteError bool

	IsSubtask bool
}

func (r *Task) PrintLn(text string) {
	text = r.Name + ": " + text
	if r.IsSubtask {
		text = "  " + text
	}
	if r.Spinner != nil {
		r.Spinner.UpdateMessage(text)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, text+"\n")
	}
}
func (r *Task) SetSuccessMessage(text string) {
	text = r.Name + ": " + text
	if r.IsSubtask {
		text = "  " + text
	}
	if r.Spinner != nil {
		r.Spinner.UpdateMessage(text)
		r.Spinner.Complete()
	}
}
func (r *Task) SetFailMessage(text string) {
	text = r.Name + ": " + text
	if r.IsSubtask {
		text = "  " + text
	}
	if r.Spinner != nil {
		r.Spinner.UpdateMessage(text)
		r.Spinner.Error()
	}
}

type TaskRunner struct {
	spinnerManager ysmrr.SpinnerManager
}

func (t *TaskRunner) GetNewSubtaskSpinner(name string) *ysmrr.Spinner {
	if t.spinnerManager == nil {
		return nil
	}
	spinner := t.spinnerManager.AddSpinner("  " + name)
	return spinner
}

func (t *TaskRunner) RunSubTasks(task Task) error {
	task.TaskRunner = t
	var errors []string
	for _, subTask := range task.SubTasks {
		subTask.Spinner = t.spinnerManager.AddSpinner(subTask.Name)
		subTaskError := subTask.Run(&subTask)
		t.updateSpinnerStatus(subTask, subTaskError)
		if subTaskError != nil {
			if !task.RunAllSubTasksDespiteError {
				return subTaskError
			} else {
				errors = append(errors, subTaskError.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func (t *TaskRunner) updateSpinnerStatus(task Task, err error) {
	if task.Spinner == nil {
		return
	}
	if task.Spinner.IsComplete() || task.Spinner.IsError() {
		return
	}
	if err == nil {
		task.Spinner.Complete()
		return
	}
	if task.ErrorAsErrorMessage {
		task.Spinner.UpdateMessage(err.Error())
	}
	task.Spinner.Error()
}

func (t *TaskRunner) RunTask(task Task) error {
	if task.Disabled {
		return nil
	}
	task.TaskRunner = t
	useSpinner := !viper.GetBool("verbose")
	if !useSpinner {
		_, _ = fmt.Fprintf(os.Stdout, "⌛ %s\n", task.Name)
	} else {
		t.spinnerManager = ysmrr.NewSpinnerManager(
			ysmrr.WithAnimation(animations.Dots),
			ysmrr.WithFrameDuration(150*time.Millisecond),
			ysmrr.WithSpinnerColor(colors.FgHiGreen),
			ysmrr.WithCompleteColor(colors.FgHiGreen),
			ysmrr.WithErrorColor(colors.FgHiRed),
		)
		t.spinnerManager.Start()
		spinner := t.spinnerManager.AddSpinner(task.Name)
		task.Spinner = spinner
	}

	err := task.Run(&task)
	if err == nil && len(task.SubTasks) > 0 {
		subTaskErrors := t.RunSubTasks(task)
		if subTaskErrors != nil {
			if !task.IgnoreSubtaskErrors {
				err = subTaskErrors
			}
		}
	}
	t.updateSpinnerStatus(task, err)

	if useSpinner {
		t.spinnerManager.Stop()
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "✗ %s\n", err.Error())
		return err
	}
	return nil
}
