package ansible

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GetAnsibleVersion returns the installed Ansible version. Returns error if not installed
func GetAnsibleVersion() (string, error) {
	cmd := exec.Command("ansible", "--version")
	var stdoutBuffer = new(bytes.Buffer)
	cmd.Stdout = stdoutBuffer
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return extractAnsibleVersion(stdoutBuffer.String())
}

// extractAnsibleVersion extracts the version from version output
func extractAnsibleVersion(versionCmdOutput string) (string, error) {
	lines := strings.Split(versionCmdOutput, "\n")
	versionLine := lines[0]

	// First line outputs e.g. "ansible 2.10.0"
	r := regexp.MustCompile(`\w+ (\d+.\d+.\d+)`)
	return r.FindStringSubmatch(versionLine)[1], nil
}

// GetCollectionDirs returns a list of Ansible collection paths from config or environment
func GetCollectionDirs() ([]string, error) {
	var customCollectionPath = viper.GetString("ansible.collection_path")
	if customCollectionPath != "" {
		absPath, err := filepath.Abs(customCollectionPath)
		if err != nil {
			return nil, err
		}
		return []string{absPath}, nil
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pathEnv = os.Getenv("COLLECTIONS_PATHS")
	var pathList []string
	if pathEnv == "" {
		pathList = []string{homeDir + "/.ansible/collections/ansible_collections"}
	} else {
		pathList = strings.Split(pathEnv, ":")
		if len(pathList) == 0 {
			pathList = []string{homeDir + "/.ansible/collections/ansible_collections"}
		}
	}
	return pathList, nil
}

// GetStackHeadCollectionLocation returns the exact path of where the StackHead collection has been installed to
func GetStackHeadCollectionLocation() (string, error) {
	collectionDirs, err := GetCollectionDirs()
	if err != nil {
		return "", err
	}
	// Look for Ansible directory
	for _, singlePath := range collectionDirs {
		var installPath = filepath.Join(singlePath, "getstackhead", "stackhead")
		if _, err := os.Stat(installPath); !os.IsNotExist(err) {
			return installPath, nil
		}
	}
	return "", fmt.Errorf("unable to find StackHead Ansible collection. Make sure to install it with \"stackhead-cli init\"")
}
