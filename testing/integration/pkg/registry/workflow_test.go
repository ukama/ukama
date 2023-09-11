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
	// w.RegisterTestCase(TC_registry_update_member)

	/* Add Network */
	w.RegisterTestCase(TC_registry_add_network)

	/* Get Network */
	w.RegisterTestCase(TC_registry_get_network)

	/* Get Networks */
	w.RegisterTestCase(TC_registry_get_networks)

	/* Add Site */
	w.RegisterTestCase(TC_registry_add_site)

	/* Get Sites */
	w.RegisterTestCase(TC_registry_get_sites)

	/* Get Site */
	w.RegisterTestCase(TC_registry_get_site)

	// /* Add Invite */
	// w.RegisterTestCase(TC_registry_add_invite)

	// /* Update Invite */
	// w.RegisterTestCase(TC_registry_update_invite)

	// /* Get Invite */
	// w.RegisterTestCase(TC_registry_get_invite)

	// /* Get Invites by org */
	// w.RegisterTestCase(TC_registry_get_invites)

	/* Add Node */
	w.RegisterTestCase(TC_registry_add_node("parent"))
	w.RegisterTestCase(TC_registry_add_node("left"))
	w.RegisterTestCase(TC_registry_add_node("right"))

	/* Update Node */
	// w.RegisterTestCase(TC_registry_update_node)

	/* Update Node State */
	w.RegisterTestCase(TC_registry_update_node_state)

	/* Add Node to site*/
	w.RegisterTestCase(TC_registry_add_node_to_site("parent"))
	w.RegisterTestCase(TC_registry_add_node_to_site("left"))
	w.RegisterTestCase(TC_registry_add_node_to_site("right"))

	/* Attach Node */
	w.RegisterTestCase(TC_registry_attach_node)

	/* Detach Node */
	w.RegisterTestCase(TC_registry_detach_node)

	/* Get node*/
	w.RegisterTestCase(TC_registry_get_node)

	/* Get nodes by site*/
	w.RegisterTestCase(TC_registry_get_nodes_by_site)

	/* Get nodes*/
	w.RegisterTestCase(TC_registry_get_nodes)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
