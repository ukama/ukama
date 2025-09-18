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

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	thttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GitClient interface {
	CreateTempDir() bool
	SetupDir() bool
	RemoveTempDir() bool
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

func NewGitClient(url, username, token, path string) *gitClient {
	return &gitClient{
		url:      url,
		username: username,
		token:    token,
		rootPath: path,
		repo:     nil,
	}
}

func (g *gitClient) RemoveTempDir() bool {
	err := os.RemoveAll(g.rootPath)
	if err != nil {
		log.Printf("remove temp dir failed: %v", err)
		return false
	}
	return true
}

func (g *gitClient) CreateTempDir() bool {
	if err := os.MkdirAll(g.rootPath, os.ModePerm); err != nil {
		log.Printf("create temp dir failed: %v", err)
		return false
	}
	return true
}

func (g *gitClient) SetupDir() bool {
	r := g.RemoveTempDir()
	c := g.CreateTempDir()
	if !r || !c {
		return false
	}
	return true
}

func (g *gitClient) CloneGitRepo(branch string) error {
	// Use 'main' as default branch if not provided
	if branch == "" {
		branch = "main"
	}

	fmt.Print(g.rootPath, g.username, g.token, g.url)
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
	g.repo = r

	return err
}

func (g *gitClient) BranchCheckout(branch string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}

	if err := w.Checkout(&branchCoOpts); err != nil {
		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		err = fetchOrigin(g.repo, mirrorRemoteBranchRefSpec, g.username, g.token)
		if err != nil {
			return err
		}

		err = w.Checkout(&branchCoOpts)
		if err != nil {
			return err
		}

		log.Printf("branch checkout success %s", branch)
	}

	return nil
}

func (g *gitClient) ReadFileJSON(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var jsonObj map[string]interface{}
	err = json.Unmarshal(content, &jsonObj)
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal json. Error %s", err.Error())
	}

	return json, nil
}

func (g *gitClient) ReadFileYML(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func fetchOrigin(repo *git.Repository, refSpecStr string, username string, token string) error {
	remotes, err := repo.Remotes()
	if err != nil {
		return err
	}

	remote, err := repo.Remote(remotes[0].Config().Name)
	if err != nil {
		return err
	}

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
			log.Println("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

func (g *gitClient) GetFilesPath(key string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(g.rootPath,
		func(_path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if strings.Contains(_path, key) {
					paths = append(paths, _path)
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
