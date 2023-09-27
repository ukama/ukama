package providers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ukama/ukama/systems/node/configurator/pkg/utils"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

type StoreProvider interface {
	GetLatestRemoteConfigs() (string, error)
	GetRemoteConfigVersion(version string) error
	GetDiff(prevSha string, curSha string, dir string) ([]string, error)
}

type gitClient struct {
	url  string
	user string
	pat  string
}

const LATEST_DIR_NAME = "/tmp/configstore/latest"
const COMMIT_DIR_NAME = "/tmp/configstore/commit"
const PERM = 0755

func (g *gitClient) GetLatestRemoteConfigs() (string, error) {

	err := utils.CreateDir(LATEST_DIR_NAME, PERM)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	r, err := git.PlainClone(LATEST_DIR_NAME, false, &git.CloneOptions{
		URL: g.url,
	})

	if err != nil {
		return "", fmt.Errorf("error cloning config store from %s : %v", g.url, err)
	}

	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("error getting head revision %s : %v", g.url, err)
	}

	hash := ref.Hash()

	return hash.String(), nil
}

func (g *gitClient) GetRemoteConfigVersion(ver string) error {

	err := utils.CreateDir(COMMIT_DIR_NAME, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	r, err := git.PlainClone(COMMIT_DIR_NAME, false, &git.CloneOptions{
		URL: g.url,
	})

	if err != nil {
		return fmt.Errorf("error cloning config store from %s : %v", g.url, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("error getting work tree: %v", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(ver),
	})
	if err != nil {
		return fmt.Errorf("error getting version %s: %v", ver, err)
	}
	return nil
}

func (g *gitClient) GetDiff(prevSha string, curSha string, dir string) ([]string, error) {

	//fmt.Println("git ")
	//CheckArgs("<revision1>")

	var hash plumbing.Hash

	//prevSha := os.Args[1] //prevSha

	// dir, err := os.Getwd()
	// CheckIfError(err)

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	// if len(os.Args) < 3 {
	// 	headRef, err := repo.Head()
	// 	CheckIfError(err)
	// 	// ... retrieving the head commit object
	// 	hash = headRef.Hash()
	// 	CheckIfError(err)
	// } else {
	// 	arg2 := os.Args[2] //optional descendent sha
	// 	hash = plumbing.NewHash(arg2)
	// }

	hash = plumbing.NewHash(curHash)
	prevHash := plumbing.NewHash(prevSha)

	prevCommit, err := repo.CommitObject(prevHash)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	commit, err := repo.CommitObject(hash)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	fmt.Println("Comparing from:" + prevCommit.Hash.String() + " to:" + commit.Hash.String())

	isAncestor, err := commit.IsAncestor(prevCommit)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	fmt.Printf("Is the prevCommit an ancestor of commit? : %v %v\n", isAncestor)

	currentTree, err := commit.Tree()
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	prevTree, err := prevCommit.Tree()
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	patch, err := currentTree.Patch(prevTree)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	fmt.Println("Got here" + strconv.Itoa(len(patch.Stats())))

	var changedFiles []string
	for _, fileStat := range patch.Stats() {
		fmt.Println(fileStat.Name)
		changedFiles = append(changedFiles, fileStat.Name)
	}

	// changes, err := currentTree.Diff(prevTree)
	// if err != nil {
	// 	log.Errorf("Error: %v", err)
	// 	return nil, err
	// }

	// for _, change := range changes {
	// 	// Ignore deleted files
	// 	action, err := change.Action()
	// 	if err != nil {
	// 		log.Errorf("Error: %v", err)
	// 		return err
	// 	}
	// 	if action == merkletrie.Delete {
	// 		//fmt.Println("Skipping delete")
	// 		continue
	// 	}

	// 	// Get list of involved files
	// 	name := getChangeName(change)
	// 	fmt.Println(name)
	// }

	return changedFiles, nil
}

func getChangeName(change *object.Change) string {
	var empty = object.ChangeEntry{}
	if change.From != empty {
		return change.From.Name
	}

	return change.To.Name
}

func NewStoreClient(url string, user string, pat string, t time.Duration) (*gitClient, error) {

	N := &gitClient{
		url:  url,
		user: user,
		pat:  pat,
	}

	return N, nil
}
