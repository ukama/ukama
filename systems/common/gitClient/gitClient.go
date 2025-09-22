/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package gitClient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	thttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	DefaultBranch  = "main"
	DefaultTimeout = 5 * time.Minute
)

type GitClient interface {
	CreateTempDir() error
	SetupDir() error
	RemoveTempDir() error
	CloneGitRepo(branch string) error
	BranchCheckout(branch string) error
	ReadFileJSON(path string) ([]byte, error)
	ReadFileYML(path string) ([]byte, error)
	GetFilesPath(key string) ([]string, error)
}

type gitClient struct {
	url      string
	username string
	token    string
	rootPath string
	repo     *git.Repository
}

func NewGitClient(url, username, token, path string) (*gitClient, error) {
	if url == "" {
		return nil, fmt.Errorf("git URL cannot be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("root path cannot be empty")
	}

	return &gitClient{
		url:      url,
		username: username,
		token:    token,
		rootPath: path,
		repo:     nil,
	}, nil
}

func (g *gitClient) RemoveTempDir() error {
	if g.rootPath == "" {
		return fmt.Errorf("root path is empty, cannot remove directory")
	}

	err := os.RemoveAll(g.rootPath)
	if err != nil {
		return fmt.Errorf("failed to remove temp directory %s: %w", g.rootPath, err)
	}

	log.Printf("Successfully removed temp directory: %s", g.rootPath)
	return nil
}

func (g *gitClient) CreateTempDir() error {
	if g.rootPath == "" {
		return fmt.Errorf("root path is empty, cannot create directory")
	}

	if err := os.MkdirAll(g.rootPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory %s: %w", g.rootPath, err)
	}

	log.Printf("Successfully created temp directory: %s", g.rootPath)
	return nil
}

func (g *gitClient) SetupDir() error {
	// Try to remove existing directory (ignore error if it doesn't exist)
	if err := g.RemoveTempDir(); err != nil {
		log.Printf("Warning: failed to remove existing directory: %v", err)
	}

	// Create new directory
	if err := g.CreateTempDir(); err != nil {
		return fmt.Errorf("failed to setup directory: %w", err)
	}

	return nil
}

func (g *gitClient) CloneGitRepo(branch string) error {
	if branch == "" {
		branch = DefaultBranch
	}

	// Validate inputs
	if g.url == "" {
		return fmt.Errorf("git URL is empty")
	}
	if g.rootPath == "" {
		return fmt.Errorf("root path is empty")
	}

	log.Printf("Cloning repository from %s, branch: %s to %s", g.url, branch, g.rootPath)

	r, err := git.PlainClone(g.rootPath, false, &git.CloneOptions{
		Auth: &thttp.BasicAuth{
			Username: g.username,
			Password: g.token,
		},
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		URL:           g.url,
		Progress:      os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository %s (branch: %s): %w", g.url, branch, err)
	}

	g.repo = r
	log.Printf("Successfully cloned repository from %s, branch: %s", g.url, branch)
	return nil
}

func (g *gitClient) BranchCheckout(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if g.repo == nil {
		return fmt.Errorf("repository is not initialized, clone repository first")
	}

	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}

	log.Printf("Attempting to checkout branch: %s", branch)

	if err := w.Checkout(&branchCoOpts); err != nil {
		log.Printf("Initial checkout failed, fetching branch from remote: %v", err)

		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		if err := fetchOrigin(g.repo, mirrorRemoteBranchRefSpec, g.username, g.token); err != nil {
			return fmt.Errorf("failed to fetch branch %s from remote: %w", branch, err)
		}

		if err := w.Checkout(&branchCoOpts); err != nil {
			return fmt.Errorf("failed to checkout branch %s after fetch: %w", branch, err)
		}
	}

	log.Printf("Successfully checked out branch: %s", branch)
	return nil
}

func (g *gitClient) ReadFileJSON(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Validate JSON without unnecessary unmarshal/marshal cycle
	var jsonObj interface{}
	if err := json.Unmarshal(content, &jsonObj); err != nil {
		return nil, fmt.Errorf("invalid JSON in file %s: %w", path, err)
	}

	// Return the original content if it's valid JSON
	return content, nil
}

func (g *gitClient) ReadFileYML(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", path, err)
	}

	return content, nil
}

func fetchOrigin(repo *git.Repository, refSpecStr string, username string, token string) error {
	if repo == nil {
		return fmt.Errorf("repository is nil")
	}
	if refSpecStr == "" {
		return fmt.Errorf("refSpec cannot be empty")
	}
	if username == "" || token == "" {
		return fmt.Errorf("authentication credentials cannot be empty")
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes: %w", err)
	}
	if len(remotes) == 0 {
		return fmt.Errorf("no remotes found in repository")
	}

	remote, err := repo.Remote(remotes[0].Config().Name)
	if err != nil {
		return fmt.Errorf("failed to get remote %s: %w", remotes[0].Config().Name, err)
	}

	log.Printf("Fetching refSpec: %s from remote: %s", refSpecStr, remote.Config().Name)

	if err = remote.Fetch(&git.FetchOptions{
		Auth: &thttp.BasicAuth{
			Username: username,
			Password: token,
		},
		RemoteURL: remote.Config().URLs[0],
		RefSpecs: []config.RefSpec{
			config.RefSpec(refSpecStr),
		},
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("References already up to date")
			return nil
		}
		return fmt.Errorf("failed to fetch from remote %s: %w", remote.Config().Name, err)
	}

	log.Printf("Successfully fetched refSpec: %s", refSpecStr)
	return nil
}

func (g *gitClient) GetFilesPath(key string) ([]string, error) {
	if key == "" {
		return nil, fmt.Errorf("search key cannot be empty")
	}
	if g.rootPath == "" {
		return nil, fmt.Errorf("root path is empty")
	}

	paths := []string{}
	err := filepath.Walk(g.rootPath,
		func(_path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking path %s: %w", _path, err)
			}
			if !info.IsDir() {
				if strings.Contains(_path, key) {
					paths = append(paths, _path)
				}
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", g.rootPath, err)
	}

	log.Printf("Found %d files matching key '%s' in %s", len(paths), key, g.rootPath)
	return paths, nil
}
