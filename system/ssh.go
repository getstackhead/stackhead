package system

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GenerateSSHKeyPair() (*rsa.PrivateKey, error) {
	bitSize := 4096

	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Debugln("GenerateSSHKeyPair#GenerateKey", err)
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		log.Debugln("GenerateSSHKeyPair#ValidateKey", err)
		return nil, err
	}
	return privateKey, nil
}

func getRemoteUser() string {
	if Context.CurrentAction == ContextActionServerSetup {
		return "root"
	}
	return "stackhead"
}

func ResolveRemoteUserIntoUid(username string) (int, error) {
	output, _, err := RemoteRun("id", RemoteRunOpts{Args: []string{"-u " + username}})
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(output.String())
}

func ResolveRemoteGroupIntoGid(groupname string) (int, error) {
	output, _, err := RemoteRun("getent", RemoteRunOpts{Args: []string{"group", groupname}})
	if err != nil {
		return -1, err
	}
	splitOutput := strings.Split(output.String(), ":")
	return strconv.Atoi(splitOutput[2])
}

type RemoteRunOpts struct {
	Args         []string
	WorkingDir   string
	Confidential bool
	Sudo         bool
	AllowFail    bool
	User         string
}

func RemoteRun(cmd string, opts RemoteRunOpts) (bytes.Buffer, bytes.Buffer, error) {
	user := getRemoteUser()
	if opts.User != "" {
		user = opts.User
	}
	remoteCmd := fmt.Sprintf("%s %s", cmd, strings.Join(opts.Args, " "))

	if opts.Confidential {
		log.Debugln(fmt.Sprintf("SSH [%s@%s]: %s", user, Context.TargetHost, cmd+" <omitted arguments>"))
	} else {
		if opts.WorkingDir != "" {
			log.Debugln(fmt.Sprintf("SSH [%s@%s:%s]: %s", user, Context.TargetHost, opts.WorkingDir, remoteCmd))
		} else {
			log.Debugln(fmt.Sprintf("SSH [%s@%s]: %s", user, Context.TargetHost, remoteCmd))
		}
	}

	var cmdArgs []string
	if user == "stackhead" {
		cmdArgs = []string{"-i", Context.Authentication.GetPrivateKeyPath()}
	}

	cmdArgs = append(cmdArgs, fmt.Sprintf("%s@%s", user, Context.TargetHost))

	if opts.Sudo {
		remoteCmd = "sudo " + remoteCmd
	}

	if opts.AllowFail {
		remoteCmd = "(" + remoteCmd + " || true)"
	}

	if opts.WorkingDir != "" {
		remoteCmd = "cd " + opts.WorkingDir + "; " + remoteCmd
	}
	cmdArgs = append(cmdArgs, remoteCmd)
	command := exec.Command("ssh", cmdArgs...)
	var out, outErr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &outErr
	err := command.Run()
	return out, outErr, err
}

func SimpleRemoteRun(cmd string, opts RemoteRunOpts) (string, error) {
	stdout, stderr, err := RemoteRun(cmd, opts)
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf(stderr.String())
		}
		return "", err
	}
	return stdout.String(), nil
}
