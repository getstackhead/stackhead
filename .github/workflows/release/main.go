package main

import (
	"log"
	"os"

	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

type HooksLogger struct {
	logger *log.Logger
}

func (t *HooksLogger) Init(m map[string]string) error {
	return nil
}

func (t *HooksLogger) Name() string {
	return "logger"
}

func (t *HooksLogger) Version() string {
	return "dev"
}

func (t *HooksLogger) Success(config *hooks.SuccessHookConfig) error {
	t.logger.Println("old version: " + config.PrevRelease.Version)
	t.logger.Println("new version: " + config.NewRelease.Version)
	t.logger.Printf("commit count: %d\n", len(config.Commits))
	return nil
}

func (t *HooksLogger) NoRelease(config *hooks.NoReleaseConfig) error {
	t.logger.Println("reason: " + config.Reason.String())
	t.logger.Println("message: " + config.Message)
	return nil
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &HooksLogger{
				logger: log.New(os.Stderr, "", 0),
			}
		},
	})
}
