package commands

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/plugins"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/stackhead"
	"github.com/getstackhead/stackhead/system"
	"github.com/getstackhead/stackhead/terraform"
)

func basicSetup() {
	// Update apt caches
}

func folderSetup() error {
	// Create StackHead root folder
	if err := xfs.CreateFolder("ssh://" + config.RootDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.RootDirectory, 1412, 1412); err != nil {
		return err
	}

	// Create certificates folder
	if err := xfs.CreateFolder("ssh://" + config.CertificatesDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.CertificatesDirectory, 1412, 1412); err != nil {
		return err
	}
	return nil
}

func setupSshKeys() error {
	// Create local StackHead folder
	localRemoteKeyDir := system.Context.Authentication.LocalAuthenticationDir
	if _, err := os.Stat(localRemoteKeyDir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(localRemoteKeyDir, os.ModeDir|os.ModePerm); err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("SSH keys used for using the \"stackhead\" user will be stored at: %s", localRemoteKeyDir))
	}

	localPrivateKeyPath := system.Context.Authentication.GetPrivateKeyPath()
	if _, err := os.Stat(localPrivateKeyPath); !os.IsExist(err) {
		// Create SSH key pair for stackhead user
		privateKey, err := system.GenerateSSHKeyPair()
		if err != nil {
			return err
		}
		// Save private key in PEM format
		err = ioutil.WriteFile(
			localPrivateKeyPath,
			pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			}),
			0600,
		)
		if err != nil {
			logger.Debugln(err)
			return err
		}
		publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
		if err != nil {
			return err
		}
		// Save public key in PEM format
		return ioutil.WriteFile(
			system.Context.Authentication.GetPublicKeyPath(),
			ssh.MarshalAuthorizedKey(publicRsaKey),
			0600,
		)
	}

	return nil
}

func userSetup() error {
	// Add stackhead group
	if _, _, err := system.RemoteRun("groupadd", "--system stackhead --gid 1412 -f"); err != nil {
		return fmt.Errorf("unable to add stackhead group")
	}

	// Add stackhead user if it does not exist
	if _, _, err := system.RemoteRun("id", "stackhead"); err != nil {
		if _, _, err := system.RemoteRun("adduser", "--system --shell /bin/sh --uid 1412 --no-create-home --home=/stackhead --gid 1412 stackhead"); err != nil {
			return fmt.Errorf("unable to add stackhead user")
		}
	}

	// Add stackhead user to www-data group
	if _, _, err := system.RemoteRun("usermod", "-a -G www-data stackhead"); err != nil {
		return fmt.Errorf("unable to add stackhead user to www-data group")
	}

	// Set includedir in sudoers
	sudoersInclude := "\n#includedir /etc/sudoers.d\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers", sudoersInclude, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to append to sudoers file")
	}
	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers"); err != nil {
		return fmt.Errorf("unable to validate sudoers file")
	}

	// Add public key to stackhead user
	publicKeyBytes, err := os.ReadFile(system.Context.Authentication.GetPublicKeyPath())
	if err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to read local stackhead public SSH key")
	}
	if err := xfs.CreateFolder("ssh:///stackhead/.ssh"); err != nil {
		return err
	}
	return xfs.WriteFile(
		"ssh:///stackhead/.ssh/authorized_keys",
		string(publicKeyBytes),
	)
}

// SetupServer is a command object for Cobra that provides the setup command
var SetupServer = &cobra.Command{
	Use:     "setup [ipv4 address]",
	Example: "setup 192.168.178.14",
	Short:   "Prepare a server for deployment",
	Long:    `setup will install all required software on a server. You are then able to deploy projects onto it.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		system.InitializeContext(args[0], system.ContextActionServerSetup, nil)

		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Setting up server at IP \"%s\"", args[0]),
			Run: func(r routines.RunningTask) error {
				var err error

				//deployedProjects, err := stackhead.GetDeployedProjects()
				//if err != nil {
				//	return err
				//}

				//- import_tasks: "../roles/stackhead_module_api/tasks_internal/load-all-modules-config.yml"
				//
				//- import_tasks: "../roles/stackhead_setup/tasks/facts-deploy.yml"
				r.PrintLn("Preparing setup.")
				basicSetup()

				r.PrintLn("Setting up SSH keys.")
				if err := setupSshKeys(); err != nil {
					fmt.Println("\nUnable to setup SSH keys. (" + err.Error() + ")")
					os.Exit(1)
				}

				r.PrintLn("Setting up users.")
				if err := userSetup(); err != nil {
					fmt.Println("\nUnable to create StackHead users. (" + err.Error() + ")")
					os.Exit(1)
				}

				r.PrintLn("Setting up folders.")
				if err := folderSetup(); err != nil {
					fmt.Println("\nUnable to create folders. (" + err.Error() + ")")
					os.Exit(1)
				}

				r.PrintLn("Setting up Terraform.")
				if err := terraform.Setup(); err != nil {
					fmt.Println("\nUnable to setup Terraform. (" + err.Error() + ")")
					os.Exit(1)
				}

				if err := stackhead.WriteVersion(); err != nil {
					fmt.Println("\nUnable to write version. (" + err.Error() + ")")
					os.Exit(1)
				}

				//- import_tasks: "../roles/stackhead_module_api/tasks_internal/setup.yml"

				return err
			},
			ErrorAsErrorMessage: true,
		})

		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Setting up StackHead Plugins at \"%s\"", args[0]),
			Run: func(r routines.RunningTask) error {
				var err error

				p, err := plugins.LoadPlugins()
				if err != nil {
					return err
				}

				r.PrintLn("Installing Terraform providers for plugins.")
				if err := terraform.BuildAndWriteProviders(p); err != nil {
					return err
				}

				r.PrintLn("Installing Terraform providers")
				if err := terraform.InstallProviders(); err != nil {
					return err
				}

				for _, plugin := range p {
					if plugin.SetupProgram != nil {
						r.PrintLn("Setup StackHead plugin " + plugin.Name)
						if err := plugin.SetupProgram.Run(nil); err != nil {
							return err
						}
					}
				}

				return err
			},
			ErrorAsErrorMessage: true,
		})
	},
}
