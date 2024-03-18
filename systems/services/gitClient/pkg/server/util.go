/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	thttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

type Company struct {
	Company       string `json:"company"`
	GitBranchName string `json:"git_branch_name"`
	Email         string `json:"email"`
}

type Environment struct {
	Production []Company `json:"production"`
	Test       []Company `json:"test"`
}

type Component struct {
	Company       string `json:"company"`
	Category      string `json:"category"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer"`
	Managed       string `json:"managed"`
	Warranty      uint32 `json:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}

func httpRequest(token string, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func RemoveTempDir(path string) {
	os.RemoveAll("temp")
}

func CreateTempDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func RemoveTempDirIfExist(path string) {
	RemoveTempDir(path)
	CreateTempDir(path)
}

func ReadRootFile(token string, url string) (*Environment, error) {
	body, err := httpRequest(token, url)
	if err != nil {
		return nil, err
	}

	var env Environment
	err = json.Unmarshal(body, &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func CloneGitRepo(url string, username string, token string, path string) (*git.Repository, error) {
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &thttp.BasicAuth{
			Username: username,
			Password: token,
		},

		SingleBranch: true,
		URL:          url,
		Progress:     os.Stdout,
	})

	return r, err
}

func BranchCheckout(r *git.Repository, branch string, username string, token string) error {
	_, err := r.Head()
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}
	if err := w.Checkout(&branchCoOpts); err != nil {
		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		err = fetchOrigin(r, mirrorRemoteBranchRefSpec, username, token)
		CheckIfError(err)

		err = w.Checkout(&branchCoOpts)
		CheckIfError(err)

		log.Infof("branch checkout success %s", branch)
	}
	CheckIfError(err)

	_, err = r.Head()
	CheckIfError(err)
	return nil
}

func fetchOrigin(repo *git.Repository, refSpecStr string, username string, token string) error {
	remotes, err := repo.Remotes()
	CheckIfError(err)

	remote, err := repo.Remote(remotes[0].Config().Name)
	CheckIfError(err)

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
			log.Info("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

func GetFilesPath(lpath string) (map[string][]string, error) {
	set := make(map[string][]string)
	set["accounting"] = []string{}
	set["components"] = []string{}
	set["contracts"] = []string{}
	err := filepath.Walk(lpath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				rp := strings.Replace(path, lpath, "", -1)
				parts := strings.Split(rp, "/")
				if len(parts) > 2 && (parts[0] == "components" || parts[0] == "contracts") {
					set[parts[0]] = append(set[parts[0]], path)
				} else if len(parts) > 1 && parts[0] == "accounting" {
					set[parts[0]] = append(set[parts[0]], path)
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return set, nil
}

func ReadFile(path string, company string) (Component, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return Component{}, err
	}
	var component Component
	err = yaml.Unmarshal(content, &component)
	if err != nil {
		log.Fatal(err)
	}
	component.Company = company

	return component, nil
}
