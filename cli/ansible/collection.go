package ansible

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

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
	return "", fmt.Errorf("unable to find StackHead Ansible collection ")
}
