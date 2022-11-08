package config

import (
	"os"
	"path"
)

var GetPluginDir = func() (string, error) {
	// Find home directory.
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// Create the file
	return path.Join(configDir, "stackhead", "plugins"), nil
}
