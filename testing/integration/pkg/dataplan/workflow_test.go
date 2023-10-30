/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package dataplan

import (
	"context"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/testing/integration/pkg/test"
)

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
}
func TestWorkflow_DataPlanSystem(t *testing.T) {

	/* Sim pool */
	w := test.NewWorkflow("dataplan_workflow_1", "Adding rates and packages")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	/* Add baserate */
	w.RegisterTestCase(TC_dp_add_baserate)

	/* Get baserate by Id */
	w.RegisterTestCase(TC_dp_get_baserate_by_id)

	/* Get rates by Period */
	w.RegisterTestCase(TC_dp_get_baserate_by_period)

	/* Get rates by Country */
	w.RegisterTestCase(TC_dp_get_baserate_by_country)

	// Add Mark ups
	w.RegisterTestCase(TC_dp_add_markup)

	/* Get Mark up */
	w.RegisterTestCase(TC_dp_get_markup)

	/* Get rate */
	// w.RegisterTestCase(TC_dp_get_rate)

	/* Add a package */
	w.RegisterTestCase(TC_dp_add_package)

	/* Get Packages */
	w.RegisterTestCase(TC_dp_get_package_for_org)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
