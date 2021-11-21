package system

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var SIZE = 1 << 15

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

func getRemoteSshAuths(user string) ([]ssh.AuthMethod, error) {
	var auths []ssh.AuthMethod
	var signers []ssh.Signer
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		localSigners, err := agent.NewClient(aconn).Signers()
		if err != nil {
			log.Errorf("Unable to load local SSH keys.")
		} else {
			signers = append(signers, localSigners...)
		}
	}

	// Load private key from local StackHead directory
	if user == "stackhead" {
		dat, err := os.ReadFile(Context.Authentication.GetPrivateKeyPath())
		if err != nil {
			return auths, err
		}
		signer, err := ssh.ParsePrivateKey(dat)
		if err != nil {
			return auths, err
		}
		signers = append(signers, signer)
	}

	auths = append(auths, ssh.PublicKeys(signers...))
	return auths, nil
}

func getRemoteClient() (*sftp.Client, error) {
	user := getRemoteUser()
	auths, err := getRemoteSshAuths(user)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", Context.TargetHost, 22), &ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // todo: can we store a hostkey to fix this?
		// https://stackoverflow.com/questions/44269142/golang-ssh-getting-must-specify-hoskeycallback-error-despite-setting-it-to-n
	})
	if err != nil {
		return nil, err
	}

	return sftp.NewClient(conn, sftp.MaxPacket(SIZE))
}

func getRemoteUser() string {
	if Context.CurrentAction == ContextActionServerSetup {
		return "root"
	}
	return "stackhead"
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
