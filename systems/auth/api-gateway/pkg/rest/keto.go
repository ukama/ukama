package rest

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type PermissionClient interface {
	CheckPermission(namespace, object, relation, subjectID string) (bool, error)
}

type permissionClient struct {
	R *resty.Client
}

func NewPermissionClient(url string, debug bool) (*permissionClient, error) {
	rc, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Failed to create PermissionClient. Error: %s", err.Error())
		return nil, err
	}

	pc := &permissionClient{
		R: rc.C,
	}

	return pc, nil
}

func (pc *permissionClient) CheckPermission(namespace, object, relation, subjectID string) (bool, error) {
	errStatus := &rest.ErrorMessage{}
	req := struct {
		Namespace  string `json:"namespace"`
		Object     string `json:"object"`
		Relation   string `json:"relation"`
		SubjectID  string `json:"subject_id"`
	}{
		Namespace:  namespace,
		Object:     object,
		Relation:   relation,
		SubjectID:  subjectID,
	}
	resp, err := pc.R.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetError(errStatus).
		Post("http://localhost:4466/check")
	if err != nil {
		logrus.Errorf("Failed to send API request to Ory keto. Error: %s", err.Error())
		return false, err
	}

	var result struct {
		Allowed bool `json:"allowed"`
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logrus.Tracef("Failed to deserialize response from Ory keto. Error message is %s", err.Error())
		return false, fmt.Errorf("Response deserialization failure: %s", err.Error())
	}
	
	return result.Allowed, nil
}