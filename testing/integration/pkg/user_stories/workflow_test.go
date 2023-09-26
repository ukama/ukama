package subscriber

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

func TestWorkflow_UserStories(t *testing.T) {
	w := test.NewWorkflow("user_stories_workflow", "User Stories Workflow")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	// w.RegisterTestCase(Story_add_org)
	w.RegisterTestCase(Story_add_user)
	w.RegisterTestCase(Story_add_network)
	w.RegisterTestCase(Story_add_node("parent"))
	w.RegisterTestCase(Story_add_node("left"))
	w.RegisterTestCase(Story_add_node("right"))
	w.RegisterTestCase(Story_add_node_to_site("parent"))
	w.RegisterTestCase(Story_add_node_to_site("left"))
	w.RegisterTestCase(Story_add_node_to_site("right"))
	w.RegisterTestCase(Story_attach_node())
	w.RegisterTestCase(Story_invite_add())           // Invite member to non-community org
	w.RegisterTestCase(Story_invite_status_update()) //Invite member to non-community org
	w.RegisterTestCase(Story_add_network_failed)
	// w.RegisterTestCase(Story_member_add()) //Add member to non-community org
	w.RegisterTestCase(Story_upload_baserate())
	w.RegisterTestCase(Story_markup())
	w.RegisterTestCase(Story_package())
	w.RegisterTestCase(Story_Simpool())
	w.RegisterTestCase(Story_Subscriber())
	w.RegisterTestCase(Story_Sim_Allocate())
	w.RegisterTestCase(Story_add_sim_package())
	w.RegisterTestCase(Story_activate_sim())
	w.RegisterTestCase(Story_active_sim_package())

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
