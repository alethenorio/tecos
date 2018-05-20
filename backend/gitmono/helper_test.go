package gitmono

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// A MockTicker allow us to control the ticks independent of actual time
// by including a Tick() method
type MockTicker struct {
	C chan time.Time
}

func NewMockTicker() *MockTicker {
	c := make(chan time.Time)
	return &MockTicker{
		C: c,
	}
}

func (m *MockTicker) Stop() {
}

func (m *MockTicker) TickerChan() <-chan time.Time {
	return m.C
}

func (m *MockTicker) Tick() {
	m.C <- time.Now()
}

// createGitRepo creates a local git repository with 1 commit and returns
// the folder path, the repo url, the repo object and a cleanup function
//
// It is up to the caller to call the cleanup function otherwise objects may be left
// in the filesystem when tests are done
func createGitRepo(dirName string, t *testing.T) (string, string, *git.Repository, func()) {
	t.Helper()

	dir, dirCleaner := testDir(dirName, t)

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Unable to create test repository: %s", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Unable to retrieve source repo working tree: %s\n", err)
	}

	_, err = wt.Add(".")
	if err != nil {
		t.Fatalf("Unable to stage path in source repo: %s\n", err)
	}

	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Unable to commit: %s\n", err)
	}

	srcURL, _ := url.Parse(filepath.ToSlash(dir))
	srcURL.Scheme = "file"

	return filepath.ToSlash(dir), srcURL.String(), repo, dirCleaner
}

func headCommit(repo *git.Repository, t *testing.T) plumbing.Hash {
	t.Helper()

	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("Unable to retrieve repo head: %s", err)
	}

	return ref.Hash()
}

// gitCommit commits a random file change in the repository in the given path
// and returns the file path and the commit hash
func gitCommit(path string, t *testing.T) (string, plumbing.Hash) {
	t.Helper()

	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatalf("Unable to open source repo: %s\n", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Unable to retrieve source repo working tree: %s\n", err)
	}

	tf, err := ioutil.TempFile(path, "testFile")
	if err != nil {
		t.Fatalf("Unable to create test file in dummy repo: %s", err)
	}

	fmt.Printf("Writing something along the lines of %s", time.Now().String())
	_, err = tf.WriteString(time.Now().String())
	if err != nil {
		t.Fatalf("Unable to write to testfile: %s", err)
	}

	fName := filepath.Base(tf.Name())

	_, err = wt.Add(fName)
	if err != nil {
		t.Fatalf("Unable to stage test file: %s", err)
	}

	hash, err := wt.Commit(fmt.Sprintf("Commited %s", fName), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		t.Fatalf("Unable to commit test file: %s", err)
	}

	return fName, hash
}

// testDir creates a test directory in the OS temp directory
// It returns the absolute path of the created directory and
// a cleanup function
//
// It is up to the caller to call the cleanup function for deleting
// the temporary directory
func testDir(dirName string, t *testing.T) (string, func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", dirName)
	if err != nil {
		t.Fatalf("Unable to create test repository: %s", err)
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}
