package system

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

type TargetType int

const (
	Pkging TargetType = iota
	Local
	Remote
)

func resolveFilePath(filePath string) (TargetType, string) {
	if strings.HasPrefix(filePath, "pkging://") {
		return Pkging, filePath[9:]
	}
	if strings.HasPrefix(filePath, "ssh://") {
		return Remote, filePath[6:]
	}
	return Local, filePath
}

func CreateFolder(folderPath string) error {
	connectionType, realPath := resolveFilePath(folderPath)
	if connectionType == Pkging {
		return fmt.Errorf("writing to Pkging is not supported and should never happen")
	}
	if connectionType == Remote {
		client, err := getRemoteClient()
		if err != nil {
			return err
		}
		defer client.Close()
		log.Debugln(fmt.Sprintf("SFTP [%s@%s]: %s", getRemoteUser(), Context.TargetHost, "MKDIRALL "+realPath))
		if err := client.MkdirAll(realPath); err != nil {
			return err
		}
		log.Debugln(fmt.Sprintf("SFTP [%s@%s]: %s", getRemoteUser(), Context.TargetHost, "CHOWN 1412:1412 "+realPath))
		if err := client.Chown(realPath, 1412, 1412); err != nil {
			return err
		}
		return nil
	}
	// Local
	return os.MkdirAll(realPath, os.ModeDir)
}

func WriteFile(filePath string, fileContent string) error {
	connectionType, realPath := resolveFilePath(filePath)
	if connectionType == Pkging {
		return fmt.Errorf("writing to Pkging is not supported and should never happen")
	}
	if connectionType == Remote {
		client, err := getRemoteClient()
		if err != nil {
			return err
		}
		defer client.Close()
		file, err := client.OpenFile(realPath, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write([]byte(fileContent)); err != nil {
			return err
		}
		if err := file.Chown(1412, 1412); err != nil {
			return err
		}
		return nil
	}
	// Local
	f, err := os.OpenFile(realPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(fileContent); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func AppendToFile(filePath string, fileContent string, onlyIfMissing bool) error {
	connectionType, realPath := resolveFilePath(filePath)
	if connectionType == Pkging {
		return fmt.Errorf("writing to Pkging is not supported and should never happen")
	}

	if onlyIfMissing {
		content, err := ReadFile(filePath)
		if err != nil {
			return err
		}
		if strings.Contains(content, fileContent) {
			return nil
		}
	}

	if connectionType == Remote {
		client, err := getRemoteClient()
		if err != nil {
			return err
		}
		defer client.Close()
		file, err := client.OpenFile(realPath, os.O_WRONLY|os.O_APPEND)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.Write([]byte(fileContent)); err != nil {
			return err
		}
		return nil
	}
	// Local
	f, err := os.OpenFile(realPath, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(fileContent); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) (string, error) {
	connectionType, realPath := resolveFilePath(filePath)
	// pull from pkging resources
	if connectionType == Pkging {
		var f pkging.File
		f, err := pkger.Open(realPath)
		if err != nil {
			return "", err
		}
		defer f.Close()
		var sl []byte
		sl, err = ioutil.ReadAll(f)
		return string(sl), nil
	}
	// pull from remote server
	if connectionType == Remote {
		client, err := getRemoteClient()
		if err != nil {
			return "", err
		}
		defer client.Close()
		file, err := client.OpenFile(realPath, os.O_RDONLY)
		if err != nil {
			return "", err
		}
		defer file.Close()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(file); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	// Local file
	dat, err := os.ReadFile(realPath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
