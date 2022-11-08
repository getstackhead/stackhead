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
	output, _, err := RemoteRun("id", "-u "+username)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(output.String())
}

func ResolveRemoteGroupIntoGid(groupname string) (int, error) {
	output, _, err := RemoteRun("getent", "group", groupname)
	if err != nil {
		return -1, err
	}
	splitOutput := strings.Split(output.String(), ":")
	return strconv.Atoi(splitOutput[2])
}

func RemoteRun(cmd string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	user := getRemoteUser()
	remoteCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	log.Debugln(fmt.Sprintf("SSH [%s@%s]: %s", getRemoteUser(), Context.TargetHost, remoteCmd))

	var cmdArgs []string
	if user == "stackhead" {
		cmdArgs = []string{"-i", Context.Authentication.GetPrivateKeyPath()}
	}

	cmdArgs = append(cmdArgs, fmt.Sprintf("%s@%s", user, Context.TargetHost))
	cmdArgs = append(cmdArgs, remoteCmd)
	command := exec.Command("ssh", cmdArgs...)
	var out, outErr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &outErr
	err := command.Run()
	return out, outErr, err
}
