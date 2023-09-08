package registry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/testing/integration/pkg/test"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)
}

func TestWorkflow_RegistrySystem(t *testing.T) {
	w := test.NewWorkflow("Registry Workflows", "Various use cases whille adding registry items")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {

		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}
	/* Get Member */
	w.RegisterTestCase(TC_registry_get_member)

	/* Get Members */
	w.RegisterTestCase(TC_registry_get_members)

	/* Update Member */
	w.RegisterTestCase(TC_registry_update_member)

	/* Add Network */
	w.RegisterTestCase(TC_registry_add_network)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
