package commandsinit

import (
	"os"
	"path"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/tools/go/vcs"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/plugins"
	"github.com/getstackhead/stackhead/routines"
)

var InstallPlugins = func() routines.Task {
	return routines.Task{
		Name:                "Installing StackHead Plugins",
		ErrorAsErrorMessage: true,
		Run: func(r routines.RunningTask) error {
			r.SetSuccessMessage("Stop message")
			r.SetFailMessage("Failed message")
			var modules = plugins.CollectPluginPaths()
			for _, modulePath := range modules {
				moduleName, moduleVersion := plugins.SplitPluginPath(modulePath)

				// Resolve repository
				rr, err := vcs.RepoRootForImportPath(moduleName, false)
				if err != nil {
					return err
				}

				pluginDir, err := config.GetPluginDir()
				if err != nil {
					r.SetFailMessage(err.Error())
					os.Exit(1)
				}
				// Create the file
				moduleSaveDir := path.Join(pluginDir, moduleName)
				err = os.MkdirAll(moduleSaveDir, 0755)
				if err != nil && !os.IsExist(err) {
					return err
				}

				// PlainOpen
				var wasCloned = false
				repo, err := git.PlainOpen(moduleSaveDir)
				if err == git.ErrRepositoryNotExists {
					r.PrintLn("Cloning repository")
					// Clone repo
					repo, err = git.PlainClone(moduleSaveDir, false, &git.CloneOptions{
						URL:           rr.Repo,
						Depth:         1,
						ReferenceName: plumbing.NewBranchReferenceName(moduleVersion),
					})
					wasCloned = true
				} else {
					r.PrintLn("Updating existing repository")
				}
				if err != nil {
					return err
				}
				if wasCloned {
					err = repo.Fetch(&git.FetchOptions{
						Depth: 1,
					})
					if err != nil {
						return err
					}
					w, err := repo.Worktree()
					if err != nil {
						return err
					}
					if err = w.Checkout(&git.CheckoutOptions{
						Branch: plumbing.NewBranchReferenceName(moduleVersion),
						Force:  true,
					}); err != nil {
						return err
					}
					if err = w.Pull(&git.PullOptions{
						Depth:         1,
						Force:         true,
						ReferenceName: plumbing.NewBranchReferenceName(moduleVersion),
					}); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
}
