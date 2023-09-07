package nucleus

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
func TestWorkflow_NucleusSystem(t *testing.T) {

	/* Sim pool */
	w := test.NewWorkflow("nucleus_workflow_1", "Adding Org and User")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	/* Add user */
	w.RegisterTestCase(TC_nucleus_add_user)

	/* Get user */
	w.RegisterTestCase(TC_nucleus_get_user)

	/* Get user by auth id */
	w.RegisterTestCase(TC_nucleus_get_user_by_auth)

	/* Whoami */
	w.RegisterTestCase(TC_nucleus_whoami)

	/* Get org */
	w.RegisterTestCase(TC_nucleus_get_org)

	/* Add user to org */
	w.RegisterTestCase(TC_nucleus_add_user_to_org)

	/* Remove user from org */
	// w.RegisterTestCase(TC_nucleus_remove_user_from_org)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
