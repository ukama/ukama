package test

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

type WorkflowExitHandlerFxn func(ctx context.Context, w *Workflow) error
type WorkflowSetupHandlerFxn func(ctx context.Context, w *Workflow) error
type WorkflowStatusHandlerFxn func(ctx context.Context, w *Workflow) (bool, error)

type ExitHandlerFxn func(ctx context.Context, t *Test) error
type SetupHandlerFxn func(ctx context.Context, t *Test) error
type StatusHanderFxn func(ctx context.Context, t *Test) (bool, error)
type TestFxn func(ctx context.Context, t *Test) error

type Workflow struct {
	Name        string
	Description string
	SetUpFxn    WorkflowSetupHandlerFxn
	CheckFxn    WorkflowStatusHandlerFxn
	ExitFxn     WorkflowExitHandlerFxn
	testSeq     []*Test
	Status      string `default:"waiting"`
	Watcher     *utils.Watcher
	Data        interface{}
}

type Test struct {
	Name        string
	Description string
	SetUpFxn    SetupHandlerFxn
	CheckFxn    StatusHanderFxn
	ExitFxn     ExitHandlerFxn
	Fxn         TestFxn
	Watcher     *utils.Watcher
	Status      string `default:"waiting"`
	Data        interface{}
	Workflow    *Workflow
}

func (t *Test) GetData() interface{} {
	return t.Data
}

func (t *Test) GetWorkflowData() interface{} {
	return t.Workflow.Data
}

func (t *Test) SaveWorkflowData(d interface{}) {
	t.Workflow.Data = d
}

func (t *Test) String() string {
	return fmt.Sprintf(` name       : %s \t\t
  desc       : %s \t\t
  status : %s
`, t.Name, t.Description, t.Status)
}

func (t *Test) Run(ctx context.Context) error {

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
			t.Status = "Fail"
			return err
		}
	} else {
		log.Errorf("Invalid test %s", t.Name)
		t.Status = "Invalid"
	}

	if t.CheckFxn != nil {
		status, err := t.CheckFxn(ctx, t)
		if err != nil {
			log.Errorf("Error while checking test %s status.", t.Name)
			t.Status = "Fail"
			return err
		}
		if status {
			t.Status = "Pass"
		} else {
			t.Status = "Fail"
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
	return fmt.Sprintf(`Workflow name : %s
description : %s
`, s.Name, s.Description)
}

func (w *Workflow) RegisterTest(t *Test) {
	w.testSeq = append(w.testSeq, t)
}

func (w *Workflow) ListTest() {
	for _, t := range w.testSeq {
		log.Infof(t.String())
	}
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
		if w.CheckFxn != nil {
			status, err := w.CheckFxn(ctx, w)
			if err != nil {
				log.Errorf("Error while checking workflow %s status.", w.Name)
				w.Status = "Fail"
			}
			if status {
				w.Status = "Pass"
			} else {
				w.Status = "Fail"
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
			w.Status = "failure"
			return err
		}
		log.Infof("Test Status: \t %s", t.String())
	}

	log.Infof("Workflow data: %+v", w.Data)
	return nil
}
