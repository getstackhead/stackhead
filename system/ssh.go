package system

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
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

func checkKnownHosts() (ssh.HostKeyCallback, error) {
	f, fErr := os.OpenFile(getHostKeyPath(), os.O_CREATE, 0600)
	if fErr != nil {
		log.Fatal(fErr)
	}
	_ = f.Close()
	return kh.New(getHostKeyPath())
}

func getClient(user string) (*ssh.Client, error) {
	var authMethod []ssh.AuthMethod
	if user == "stackhead" {
		authMethod = append(authMethod, publicKey(Context.Authentication.GetPrivateKeyPath()))
	} else {
		authMethod = append(authMethod, publicKey(path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")))
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: authMethod,
		HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
			knownHosts, err := checkKnownHosts()
			if err != nil {
				return err
			}
			var keyErr *kh.KeyError
			hErr := knownHosts(host, remote, pubKey)
			if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
				// Reference: https://www.godoc.org/golang.org/x/crypto/ssh/knownhosts#KeyError
				// if keyErr.Want slice is empty then host is unknown, if keyErr.Want is not empty
				// and if host is known then there is key mismatch the connection is then rejected.
				//log.Printf("WARNING: The received key is not a key of %s, either a MiTM attack or %s has reconfigured the host pub key.", host, host)
				return keyErr
			} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
				// host key not found in known_hosts then give a warning and continue to connect.
				//log.Printf("WARNING: %s is not trusted, adding key to known_hosts file.", host)
				return addHostKey(remote, pubKey)
			}
			//log.Printf("Pub key exists for %s.", host)
			return nil
		}),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", Context.TargetHost, "22"), config)
}

func getHostKeyPath() string {
	return path.Join(os.Getenv("HOME"), ".ssh", "stackhead_known_hosts")
}

func addHostKey(remote net.Addr, pubKey ssh.PublicKey) error {
	// add host key if host is not found in known_hosts, error object is return, if nil then connection proceeds,
	// if not nil then connection stops.
	f, fErr := os.OpenFile(getHostKeyPath(), os.O_APPEND|os.O_WRONLY, 0600)
	if fErr != nil {
		return fErr
	}
	defer f.Close()

	knownHosts := kh.Normalize(remote.String())
	_, fileErr := f.WriteString(kh.Line([]string{knownHosts}, pubKey))
	return fileErr
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

	if opts.Sudo {
		remoteCmd = "sudo " + remoteCmd
	}

	if opts.AllowFail {
		remoteCmd = "(" + remoteCmd + " || true)"
	}

	if opts.WorkingDir != "" {
		remoteCmd = "cd " + opts.WorkingDir + "; " + remoteCmd
	}

	var out, outErr bytes.Buffer
	client, err := getClient(user)
	if err != nil {
		return out, outErr, err
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		log.Fatal("unable to create SSH session: ", err)
	}
	defer ss.Close()

	ss.Stdout = &out
	ss.Stderr = &outErr
	err = ss.Run(remoteCmd)
	return out, outErr, err
}

func publicKey(path string) ssh.AuthMethod {
	key, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
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
