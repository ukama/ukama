package subscriber

import (
	"context"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	"github.com/ukama/ukama/testing/integration/pkg/dataplan"
	"github.com/ukama/ukama/testing/integration/pkg/test"
)

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
}

func TestWorkflow_SubscriberSystem(t *testing.T) {

	/* Sim pool */
	w := test.NewWorkflow("susbcriber_workflow_1", "Adding sims to sim pool")

	w.SetUpFxn = func(ctx context.Context, w *test.Workflow) error {

		log.Tracef("Initilizing Data for %s.", w.String())
		d := InitializeData()

		/* Initialize registry *
		/* add user */
		uresp, err := d.Reg.AddUser(d.reqAddUserRequest)
		if err != nil {
			return nil
		} else {
			d.UserId = uresp.User.Uuid
		}

		/* adding  */
		d.reqAddOrgRequest.Owner = d.UserId
		resp, err := d.Reg.AddOrg(d.reqAddOrgRequest)
		if err != nil {
			return nil
		} else {
			d.OrgId = resp.Org.Id
			d.OrgName = resp.Org.Name
		}

		/* adding network */
		nresp, err := d.Reg.AddNetwork(rapi.AddNetworkRequest{
			OrgName: resp.Org.Name,
			NetName: resp.Org.Name + "-net",
		})
		if err != nil {
			return nil
		} else {
			d.NetworkId = nresp.Network.Id
			d.NetworkName = nresp.Network.Name
		}

		/* Initialize the data-plan system */
		err = func(org string, owner string) error {
			dp := test.NewWorkflow("dataplan_config_for_subscriber", "Adding data pan for subscriber")

			dp.SetUpFxn = func(ctx context.Context, dp *test.Workflow) error {
				log.Tracef("Initilizing Data for %s.", dp.String())
				dp.Data = dataplan.InitializeData(&org, &owner)

				log.Tracef("Workflow Data : %+v", dp.Data)
				return nil
			}

			/* Add baserate */
			dp.RegisterTestCase(dataplan.TC_dp_add_baserate)

			// Add Mark ups
			dp.RegisterTestCase(dataplan.TC_dp_add_markup)

			/* Get rate request */
			dp.RegisterTestCase(dataplan.TC_dp_get_rate)

			/* Add a package */
			dp.RegisterTestCase(dataplan.TC_dp_add_package)

			/* Get Packages */
			dp.RegisterTestCase(dataplan.TC_dp_get_package_for_org)

			dp.ExitFxn = func(ctx context.Context, wf *test.Workflow) error {
				data := dp.GetData().(*dataplan.InitData)
				d.PackageId = data.PackageId

				return nil
			}

			/* Run */
			err := dp.Run(t, context.Background())
			assert.NoError(t, err)

			return err

		}(d.OrgId, d.UserId)

		w.Data = d
		log.Tracef("Workflow Data : %+v", w.Data)
		return err
	}

	w.RegisterTestCase(TC_simpool_upload)

	w.RegisterTestCase(TC_simpool_get_sim)

	w.RegisterTestCase(TC_simpool_get_stats)

	/* Add subscriber */
	w.RegisterTestCase(TC_registry_add_subscriber)

	/* Allocate Sim */
	w.RegisterTestCase(TC_manager_allocate_sim)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
