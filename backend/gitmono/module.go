package gitmono

import (
	"fmt"
	"os"
	"time"

	"github.com/ByteFlinger/tecos/backend"
	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

const (
	clonePath = "./gitRepo"
)

type Config struct {
	// The git repository url to be cloned and used for retrieving modules
	RepoURL string
	// The path in the git repository to retrieve modules from.
	// Use a relative path from the repo root. Wildcards can be used to include
	// several sub folders under the same folder
	ModulePath string
}

type GitMono struct {
	repoURL    string
	modulePath string
	quitChan   chan struct{}
	doneChan   chan struct{}
	repo       *git.Repository
}

func New(config *Config) (*GitMono, error) {

	// Interval to perform pull on the git repository
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	done := make(chan struct{})

	repo, err := cloneRepo(clonePath, config.RepoURL)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to clone git repository")
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to retrieve git repo working tree")
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				err := w.Pull(&git.PullOptions{})
				if err != nil && err != git.NoErrAlreadyUpToDate {
					fmt.Fprintf(os.Stderr, `Error: Unable to pull from remote repository - "%s"`, err)
				}
			case <-quit:
				ticker.Stop()
				os.RemoveAll(clonePath)
				close(done)
				return
			}
		}
	}()

	return &GitMono{
		repoURL:    config.RepoURL,
		modulePath: config.ModulePath,
		quitChan:   quit,
		doneChan:   done,
		repo:       repo,
	}, nil
}

func (d *GitMono) ListModules() []backend.ModuleData {

	return []backend.ModuleData{}

}

func (d *GitMono) Cleanup() {
	close(d.quitChan)
	<-d.doneChan
}

func cloneRepo(path, url string) (*git.Repository, error) {
	err := os.RemoveAll(path)

	if err != nil {
		return nil, errors.WithMessage(err, "Unable to remove old clone folder")
	}

	return git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}
