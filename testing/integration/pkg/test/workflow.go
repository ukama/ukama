package test

import (
	"context"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

type Workflow struct {
	Name        string
	Description string
	SetUpFxn    WorkflowSetupHandlerFxn
	StateFxn    WorkflowStateCheckHandlerFxn
	ExitFxn     WorkflowExitHandlerFxn
	testSeq     []*TestCase
	State       StateType `default:"0"`
	Watcher     *utils.Watcher
	Data        interface{}
	Count       int64 `default:"0"`
	Pass        int64 `default:"0"`
	Fail        int64 `default:"0"`
	Untested    int64 `default:"0"`
}

type TestCase struct {
	Name        string
	Description string
	SetUpFxn    SetupHandlerFxn
	StateFxn    StateHanderFxn
	ExitFxn     ExitHandlerFxn
	Fxn         TestFxn
	Watcher     *utils.Watcher
	State       StateType `default:"0"`
	Data        interface{}
	Workflow    *Workflow
}

func (t *TestCase) GetData() interface{} {
	return t.Data
}

func (t *TestCase) GetWorkflowData() interface{} {
	return t.Workflow.Data
}

func (t *TestCase) SaveWorkflowData(d interface{}) {
	t.Workflow.Data = d
}

func (t *TestCase) String() string {
	return fmt.Sprintf("Test State:: Name: %s Desc: %s Status: %s", t.Name, t.Description, t.State.String())
}

func (t *TestCase) Run(test *testing.T, ctx context.Context) error {

	log.Debugf("Starting setup for %s", t.Name)

	if t.SetUpFxn != nil {
		err := t.SetUpFxn(ctx, t)
		if assert.NoError(test, err) {
			t.State = StateTypeUnderTest
		}
	}

	if t.Fxn != nil {
		err := t.Fxn(ctx, t)
		if assert.NoError(test, err) {
			t.State = StateTypeTested
		}

	} else {
		log.Errorf("Invalid test %s", t.Name)
		t.State = StateTypeInvalid
	}

	if t.StateFxn != nil {
		status, err := t.StateFxn(ctx, t)
		if assert.NoError(test, err) {

			if assert.EqualValues(test, true, status) {
				t.State = StateTypePass
			} else {
				t.State = StateTypeFail
			}

		} else {
			t.State = StateTypeFail
		}
	}

	if t.ExitFxn != nil {
		err := t.ExitFxn(ctx, t)
		if err != nil {
			log.Errorf("Error while doing clean up  after test %s.", t.Name)
			return err
		}
	}

	return nil

}

func NewWorkflow(name, desc string) *Workflow {

	return &Workflow{
		Name:        name,
		Description: desc,
	}
}

func (w *Workflow) GetData() interface{} {
	return w.Data
}

func (w *Workflow) SaveData(d interface{}) {
	w.Data = d
}

func (w *Workflow) String() string {
	return fmt.Sprintf("Workflow: name: %s Description: %s", w.Name, w.Description)
}

func (w *Workflow) Info() string {
	return fmt.Sprintf("Workflow Name: %s Description: %s Count: %d Pass: %d Fail: %d Untested: %d", w.Name, w.Description, w.Count, w.Pass, w.Untested, w.Untested)
}

func (w *Workflow) RegisterTestCase(t *TestCase) {
	w.testSeq = append(w.testSeq, t)
	w.Count++
}

func (w *Workflow) ListTestCase() {
	for _, t := range w.testSeq {
		log.WithFields(log.Fields{
			"Name":        t.Name,
			"Description": t.Description,
			"State":       t.State,
		}).Info("Test Case: ")
	}
}

func (w *Workflow) Status() {
	log.WithFields(log.Fields{
		"Name":     w.Name,
		"Count":    w.Count,
		"Pass":     w.Pass,
		"Fail":     w.Fail,
		"Untested": w.Untested,
	}).Info(" Workflow Results")
	w.ListTestCase()
}

func (w *Workflow) Run(test *testing.T, ctx context.Context) error {

	if w.SetUpFxn != nil {
		log.Tracef("Starting setup for workflow %s", w.Name)

		err := w.SetUpFxn(ctx, w)
		if assert.NoError(test, err) {
			w.State = StateTypeUnderTest
		}
	}

	defer func() {
		if w.StateFxn != nil {
			status, err := w.StateFxn(ctx, w)
			if assert.NoError(test, err) {

				if assert.EqualValues(test, true, status) {
					w.State = StateTypePass
				} else {
					w.State = StateTypeFail
				}

			} else {
				log.Errorf("Error while checking workflow %s status.", w.Name)
				w.State = StateTypeFail
			}
		}

		if w.ExitFxn != nil {
			err := w.ExitFxn(ctx, w)
			if err != nil {
				log.Errorf("Error while doing clean up  after testing workflow %s.", w.Name)
			}
		}
	}()

	for _, tc := range w.testSeq {
		err := tc.Run(test, ctx)
		if assert.NoError(test, err) {
			w.State = StateTypeFail
		}

		w.stats(tc.State)
		log.Debugf("Test Status: %s", tc.String())
	}

	log.Tracef("Workflow data: %+v", w.Data)
	return nil
}

func (w *Workflow) stats(status StateType) {
	switch status {
	case StateTypeFail:
		w.Fail++
	case StateTypePass:
		w.Pass++
	default:
		w.Untested++
	}
}
