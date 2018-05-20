package gitmono

import (
	"fmt"
	"os"
	"time"

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

	// How often to check for new commits
	PullInterval time.Duration

	// Internal variables
	ticker Ticker
}

type GitMono struct {
	repoURL    string
	modulePath string
	quitChan   chan struct{}
	doneChan   chan struct{}
	repo       *git.Repository
}

// New returns a new gitmono.GitMono given a configuration
// Upon creation, New will clone the given repository URL in
// the current working directoryu and perform a periodic pull
// on the master branch in the given time interval
func New(config *Config) (*GitMono, error) {

	if config.PullInterval == 0 {
		config.PullInterval = 5 * time.Minute
	}

	if config.ticker == nil {
		config.ticker = NewTicker(config.PullInterval)
	}

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
		run := true
		for run {
			select {
			case <-config.ticker.TickerChan():
				err := w.Pull(&git.PullOptions{})
				if err != nil && err != git.NoErrAlreadyUpToDate {
					fmt.Fprintf(os.Stderr, "Error: Unable to pull from remote repository - \"%s\"\n", err)
				}
				fmt.Println("Pulled")
			case <-quit:
				config.ticker.Stop()
				os.RemoveAll(clonePath)
				run = false
				close(done)
				break
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
