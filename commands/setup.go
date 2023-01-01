package commands

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gookit/event"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
	"github.com/getstackhead/stackhead/terraform"
)

func basicSetup() error {
	// Install packages to allow apt to use a repository over HTTPS
	if err := system.UpdatePackageList(system.PackageVendorApt); err != nil {
		return err
	}
	if err := system.InstallPackage([]system.Package{
		{Name: "apt-transport-https", Vendor: system.PackageVendorApt},
		{Name: "ca-certificates", Vendor: system.PackageVendorApt},
		{Name: "curl", Vendor: system.PackageVendorApt},
		{Name: "debian-keyring", Vendor: system.PackageVendorApt},
		{Name: "debian-archive-keyring", Vendor: system.PackageVendorApt},
		{Name: "gnupg-agent", Vendor: system.PackageVendorApt},
		{Name: "software-properties-common", Vendor: system.PackageVendorApt},
	}); err != nil {
		return err
	}

	// Update apt caches
	return nil
}

func folderSetup() error {
	event.MustFire("setup.folders.pre-install", event.M{})
	// Create StackHead root folder
	if err := xfs.CreateFolder("ssh://" + config.RootDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.RootDirectory, 1412, 1412); err != nil {
		return err
	}
	event.MustFire("setup.folders.post-install", event.M{})
	return nil
}

func setupSshKeys() error {
	event.MustFire("setup.ssh.pre-install", event.M{})
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
		err = ioutil.WriteFile(
			system.Context.Authentication.GetPublicKeyPath(),
			ssh.MarshalAuthorizedKey(publicRsaKey),
			0600,
		)
		if err != nil {
			return err
		}
	}
	event.MustFire("setup.ssh.post-install", event.M{})

	return nil
}

func userSetup() error {
	event.MustFire("setup.users.pre-install", event.M{})
	// Add stackhead group
	if _, _, err := system.RemoteRun("groupadd", system.RemoteRunOpts{Args: []string{"--system stackhead --gid 1412 -f"}}); err != nil {
		return fmt.Errorf("unable to add stackhead group")
	}

	// Add stackhead user if it does not exist
	if _, _, err := system.RemoteRun("id", system.RemoteRunOpts{Args: []string{"stackhead"}}); err != nil {
		if _, _, err := system.RemoteRun("adduser", system.RemoteRunOpts{Args: []string{"--system --shell /bin/sh --uid 1412 --no-create-home --home=/stackhead --gid 1412 stackhead"}}); err != nil {
			return fmt.Errorf("unable to add stackhead user")
		}
	}

	// Set includedir in sudoers
	content, err := xfs.ReadFile("ssh:///etc/sudoers")
	if err != nil {
		logger.Errorln(err)
		return fmt.Errorf("unable to read sudoers file")
	}
	if !strings.Contains(content, "includedir /etc/sudoers.d\n") {
		sudoersInclude := "\nincludedir /etc/sudoers.d\n"
		if err := xfs.AppendToFile("ssh:///etc/sudoers", sudoersInclude, true); err != nil {
			logger.Errorln(err)
			return fmt.Errorf("unable to append to sudoers file")
		}
	}

	// Create empty sudoers file for additional permissions of stackhead user
	if err := xfs.WriteFile("ssh:///etc/sudoers.d/stackhead", ""); err != nil {
		return fmt.Errorf("unable to create empty stackhead sudoers file")
	}

	// todo: API to add entries to NOPASS_CMNDS
	permissions := "\nCmnd_Alias STACKHEAD_NOPASS_CMNDS = /bin/chmod, /bin/chown\n%stackhead ALL= NOPASSWD: STACKHEAD_NOPASS_CMNDS\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissions, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add chown permissions for stackhead user")
	}

	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers", system.RemoteRunOpts{}); err != nil {
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
	if err := xfs.WriteFile(
		"ssh:///stackhead/.ssh/authorized_keys",
		string(publicKeyBytes),
	); err != nil {
		return err
	}
	event.MustFire("setup.users.post-install", event.M{})
	return nil
}

// SetupServer is a command object for Cobra that provides the setup command
var SetupServer = &cobra.Command{
	Use:     "setup [ipv4 address]",
	Example: "setup 192.168.178.14",
	Short:   "Prepare a server for deployment",
	Long:    `setup will install all required software on a server. You are then able to deploy projects onto it.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		PrepareContext(args[0], system.ContextActionServerSetup, nil)

		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Setting up server at IP \"%s\"", args[0]),
			Run: func(r routines.RunningTask) error {
				var err error

				// Init modules
				for _, module := range system.Context.GetModulesInOrder() {
					moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
					module.Init(moduleSettings)
				}

				//deployedProjects, err := stackhead.GetDeployedProjects()
				//if err != nil {
				//	return err
				//}

				//- import_tasks: "../roles/stackhead_module_api/tasks_internal/load-all-modules-config.yml"
				//
				//- import_tasks: "../roles/stackhead_setup/tasks/facts-deploy.yml"
				r.PrintLn("Preparing setup.")

				if err := basicSetup(); err != nil {
					fmt.Println("\nUnable to prepare setup. (" + err.Error() + ")")
					os.Exit(1)
				}

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

				if err := system.WriteVersion(); err != nil {
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

				r.PrintLn("Installing Terraform providers for plugins.")
				if err := terraform.BuildAndWriteProviders(); err != nil {
					return err
				}

				r.PrintLn("Installing Terraform providers")
				if err := terraform.InstallProviders(); err != nil {
					return err
				}

				modules := system.Context.GetModulesInOrder()
				event.MustFire("setup.modules.pre-install", event.M{"modules": modules})
				for _, module := range modules {
					if module.GetConfig().Type == "plugin" || module.GetConfig().Type == "dns" {
						continue
					}
					event.MustFire(
						"setup.modules.pre-install-module."+module.GetConfig().Type+"."+module.GetConfig().Name,
						event.M{"module": module},
					)
					moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
					if err := module.Install(moduleSettings); err != nil {
						return err
					}
					event.MustFire(
						"setup.modules.post-install-module."+module.GetConfig().Type+"."+module.GetConfig().Name,
						event.M{"module": module},
					)
				}
				event.MustFire("setup.modules.post-install", event.M{"modules": modules})

				return err
			},
			ErrorAsErrorMessage: true,
		})
	},
}
