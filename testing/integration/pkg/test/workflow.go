package test

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
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
	return fmt.Sprintf(" Test State:: \n \t name		: %s \n \t desc		: %s \n \t status	: %s \n", t.Name, t.Description, t.State.String())
}

func (t *TestCase) Run(ctx context.Context) error {

	log.Info("Starting setup for %s", t.String())

	if t.SetUpFxn != nil {
		err := t.SetUpFxn(ctx, t)
		if err != nil {
			log.Errorf("Error while doing test setup for %s.", t.Name)
			return err
		}
	}

	if t.Fxn != nil {
		err := t.Fxn(ctx, t)
		if err != nil {
			log.Errorf("Error while executing test %s.", t.Name)
			t.State = StateTypeFail
			return err
		}
	} else {
		log.Errorf("Invalid test %s", t.Name)
		t.State = StateTypeInvalid
	}

	if t.StateFxn != nil {
		status, err := t.StateFxn(ctx, t)
		if err != nil {
			log.Errorf("Error while checking test %s status.", t.Name)
			t.State = StateTypeFail
			return err
		}
		if status {
			t.State = StateTypePass
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

func (s *Workflow) String() string {
	return fmt.Sprintf("Workflow: \n\n name 			: %s \n\t description 	: %s \n\t", s.Name, s.Description)
}

func (w *Workflow) RegisterTestCase(t *TestCase) {
	w.testSeq = append(w.testSeq, t)
}

func (w *Workflow) ListTestCase() {
	for _, t := range w.testSeq {
		log.Infof(t.String())
	}
}

func (w *Workflow) Status() {
	log.Infof(w.String())
	w.ListTestCase()
}

func (w *Workflow) Run(ctx context.Context) error {
	log.Infof("Starting workflow %s", w.String())
	if w.SetUpFxn != nil {
		log.Info("Starting setup for workflow %s", w.Name)

		err := w.SetUpFxn(ctx, w)
		if err != nil {
			log.Errorf("Error while doing workflow setup for %s.", w.Name)
			return err
		}
	}

	defer func() {
		if w.StateFxn != nil {
			status, err := w.StateFxn(ctx, w)
			if err != nil {
				log.Errorf("Error while checking workflow %s status.", w.Name)
				w.State = StateTypeFail
			}
			if status {
				w.State = StateTypePass
			} else {
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

	for _, t := range w.testSeq {
		err := t.Run(ctx)
		if err != nil {
			log.Errorf("Error while running test for workflow %s test %s.", w.Name, t.Name)
			w.State = StateTypeFail
			return err
		}
		log.Infof("Test Status: \t %s", t.String())
	}

	log.Infof("Workflow data: %+v", w.Data)
	return nil
}
