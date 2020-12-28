package main

import (
	"io/ioutil"
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

func (t *HooksLogger) writeVersionFile(version string, filePath string) error {
	return ioutil.WriteFile(filePath, []byte(version), 0644)
}

func (t *HooksLogger) Success(config *hooks.SuccessHookConfig) error {
	// adjust and copy version file ansible/VERSION /VERSION
	// copy schemas from schemas/**.json to ansible/schemas/

	// Write VERSION file
	if err := t.writeVersionFile(config.NewRelease.Version, "VERSION"); err != nil {
		return err
	}

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
	// var logger = &HooksLogger{
	// 	logger: log.New(os.Stderr, "", 0),
	// }
	// logger.Success(&hooks.SuccessHookConfig{
	// 	PrevRelease: &semrel.Release{
	// 		SHA:         "abc",
	// 		Version:     "1.0.0",
	// 		Annotations: nil,
	// 	},
	// 	NewRelease:  &semrel.Release{
	// 		SHA:         "def",
	// 		Version:     "2.0.0",
	// 		Annotations: nil,
	// 	},
	// })
	// panic("test")

	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &HooksLogger{
				logger: log.New(os.Stderr, "", 0),
			}
		},
	})
}
