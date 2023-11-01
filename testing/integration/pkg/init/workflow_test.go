/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package init

import (
	"testing"

	"context"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/testing/integration/pkg/test"
)

func TestWorkflow_InitSystem(t *testing.T) {

	w := test.NewWorkflow("init_workflow_1", "Adding a system and getting its credentials")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Debugf("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	w.RegisterTestCase(TC_init_add_org)

	w.RegisterTestCase(TC_init_add_system)

	w.RegisterTestCase(TC_init_add_node)

	w.RegisterTestCase(TC_init_bootstrap_node)

	w.RegisterTestCase(TC_init_get_system)

	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
