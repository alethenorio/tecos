package gitmono

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, rURL, _, rCleaner := createGitRepo("dummy", t)
	defer rCleaner()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"No URL", &Config{}, true},
		{"Incorrect URL", &Config{RepoURL: "https://non-existant-url"}, true},
		{"Correct URL", &Config{RepoURL: rURL}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := New(tt.config)
			close(g.quitChan)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPullInterval(t *testing.T) {

	repoDir, rURL, repo, rCleaner := createGitRepo("dummy", t)
	defer rCleaner()

	ticker := NewMockTicker()

	config := &Config{
		RepoURL: rURL,
		ticker:  ticker,
	}

	gm, err := New(config)
	defer close(gm.quitChan)

	if err != nil {
		t.Fatalf("Pull on New(): %s", err)
	}

	srcHash := headCommit(repo, t)
	dstHash := headCommit(gm.repo, t)

	fmt.Println(srcHash)

	if srcHash.String() != dstHash.String() {
		t.Fatalf("New() head commit %s does not match source repo head commit %s", dstHash.String(), srcHash.String())
	}

	_, newSrcHash := gitCommit(repoDir, t)

	ticker.Tick()
	// Wait for the pull to be finished
	// TODO(ealethe): Maybe find a more stable way to test this. Relying on a random wait time is not very good
	time.Sleep(1 * time.Second)

	newDstHash := headCommit(gm.repo, t)

	if newSrcHash.String() != newDstHash.String() {
		t.Fatalf("New() head commit %s does not match source repo head commit %s after pull", newDstHash.String(), newSrcHash.String())
	}

	//Let's try with 2 commits
	_, _ = gitCommit(repoDir, t)
	_, newSrcHash = gitCommit(repoDir, t)

	ticker.Tick()
	// Wait for the pull to be finished
	// TODO(ealethe): Maybe find a more stable way to test this. Relying on a random wait time is not very good
	time.Sleep(1 * time.Second)

	newDstHash = headCommit(gm.repo, t)

	if newSrcHash.String() != newDstHash.String() {
		t.Fatalf("New() head commit %s does not match source repo head commit %s after pull", newDstHash.String(), newSrcHash.String())
	}

}

func Test_cloneRepo(t *testing.T) {
	_, rURL, _, rCleaner := createGitRepo("dummy", t)
	defer rCleaner()

	testDir, dCleaner := testDir("test", t)
	defer dCleaner()

	// Let's run cloneRepo twice and make sure we get no errors
	_, err := cloneRepo(testDir, rURL)
	if err != nil {
		t.Errorf("cloneRepo() error = %v", err)
		return
	}

	_, err = cloneRepo(testDir, rURL)
	if err != nil {
		t.Errorf("cloneRepo() error = %v", err)
		return
	}

}
