package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/wI2L/jsondiff"
)

func JsonDiff(srcFile string, targetFile string) ([]string, bool, error) {
	change := false
	source, err := os.ReadFile(srcFile)
	if err != nil {
		log.Errorf("error reading json file %s: %v", srcFile, err)
		return nil, change, err
	}

	target, err := os.ReadFile(targetFile)
	if err != nil {
		log.Errorf("error reading json file %s: %v", targetFile, err)
		return nil, change, err
	}

	patch, err := jsondiff.CompareJSON(source, target)
	if err != nil {
		log.Errorf("error comparing json file %s and %s: %v", srcFile, targetFile, err)
		return nil, change, err
	}
	var changedValues []string
	for _, op := range patch {
		fmt.Printf("%s\n", op)
		changedValues = append(changedValues, op.String())
		change = true
	}

	return changedValues, change, nil
}
