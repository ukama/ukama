package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/wI2L/jsondiff"
)

func JsonDiff(srcFile string, targetFile string) ([]string, bool, int, error) {
	change := false
	state := 0
	source, err := os.ReadFile(srcFile)
	if err != nil {
		log.Warningf("File may be added in latest commit. Error reading json file %s.  Error %+v", srcFile, err)
		source = []byte("{}")
		state = 1 /* Added */
	}

	target, err := os.ReadFile(targetFile)
	if err != nil {
		log.Warningf("File may be deleted in latest commit. Error reading json file %s: %v", targetFile, err)
		state = 2 /* State */
		return nil, change, state, nil
	}

	patch, err := jsondiff.CompareJSON(source, target)
	if err != nil {
		log.Errorf("error comparing json file %s and %s: %v", srcFile, targetFile, err)
		return nil, change, state, err
	}

	var changedValues []string
	for _, op := range patch {
		fmt.Printf("%s\n", op)
		changedValues = append(changedValues, op.String())
		change = true
		state = 3 /* Updated */
	}

	return changedValues, change, state, nil
}
