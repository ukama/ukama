package test

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
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
	data        interface{}
}

type Test struct {
	Name        string
	Description string
	SetUpFxn    SetupHandlerFxn
	CheckFxn    StatusHanderFxn
	ExitFxn     ExitHandlerFxn
	Fxn         TestFxn
	Status      string `default:"waiting"`
	data        interface{}
}

func (t *Test) TestData() interface{} {
	return t.data
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
		log.Errorf("Inavlid test %s", t.Name)
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

	log.Info("Completed test %s", t.String())
	return nil

}

func NewWorkflow(name, desc string) *Workflow {

	return &Workflow{
		Name:        name,
		Description: desc,
	}
}

func (w *Workflow) TestData() interface{} {
	return w.data
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
		t.String()
	}
}

func (w *Workflow) Run(ctx context.Context) error {
	log.Info("Starting workflow %s", w.String())
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
		}
	}
	log.Info("Workflow %s completed", w.String())
	return nil
}
