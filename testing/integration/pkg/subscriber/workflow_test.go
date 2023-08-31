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

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {

		log.Tracef("Initilizing Data for %s.", w.String())
		d := InitializeData()

		/* Initialize registry */
		err := func() error {
			/* add user */
			uresp, err := d.Reg.AddUser(d.reqAddUserRequest)
			if err != nil {
				return err
			} else {
				d.UserId = uresp.User.Id
			}

			/* adding  */
			d.reqAddOrgRequest.Owner = d.UserId
			resp, err := d.Reg.AddOrg(d.reqAddOrgRequest)
			if err != nil {
				return err
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
				return err
			} else {
				d.NetworkId = nresp.Network.Id
				d.NetworkName = nresp.Network.Name
			}

			return nil
		}()
		if err != nil {
			log.Errorf("Initializing registry system dependencies failed. Error %v", err)
			return err
		}

		/* Initialize the data-plan system */
		err = func(org string, owner string) error {
			dp := test.NewWorkflow("dataplan_config_for_subscriber", "Adding data pan for subscriber")

			dp.SetUpFxn = func(t *testing.T, ctx context.Context, dp *test.Workflow) error {
				log.Tracef("Initializing Data for %s.", dp.String())
				dp.Data = dataplan.InitializeData(&org, &owner)

				log.Tracef("Workflow Data : %+v", dp.Data)
				return nil
			}

			/* Add baserate */
			dp.RegisterTestCase(dataplan.TC_dp_add_baserate)

			/* Add Mark ups */
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
			if err != nil {
				log.Errorf("Initializing data plan dependencies failed. Error %v", err)
				return err
			}

			d.wfDataPlan = dp
			return err

		}(d.OrgId, d.UserId)

		w.Data = d
		log.Tracef("Workflow Data for DataPlan : %+v", w.Data)
		return err
	}

	w.RegisterTestCase(TC_simpool_upload)

	w.RegisterTestCase(TC_simpool_get_sim)

	w.RegisterTestCase(TC_simpool_get_stats)

	/* Add subscriber */
	w.RegisterTestCase(TC_registry_add_subscriber)

	/* Get susbscriber */
	w.RegisterTestCase(TC_registry_get_subscriber)

	/* Allocate Sim */
	w.RegisterTestCase(TC_manager_allocate_sim)

	/* Get package for sim */
	w.RegisterTestCase(TC_manager_get_package_for_sim)

	/* Get sim: Have to call this so that test agent add a sim record*/
	w.RegisterTestCase(TC_manager_get_sim)

	/* Activate Sim */
	w.RegisterTestCase(TC_manager_activate_sim)

	/* Get Sim by subscriber */
	w.RegisterTestCase(TC_manager_get_sim_by_subscriber)

	/* Activate Package for sim */
	w.RegisterTestCase(TC_manager_set_active_package_for_sim)

	/* Check active package for sim */
	w.RegisterTestCase(TC_manager_check_active_package_for_sim)

	/* Add one more package to sim */
	w.RegisterTestCase(TC_manager_add_extra_package_to_sim)

	/* Get package for sim */
	w.RegisterTestCase(TC_manager_get_multiple_package_for_sim)

	/* Change Activate Package for sim to Additonal package */
	w.RegisterTestCase(TC_manager_set_active_package_for_sim)

	/* Check active package for sim */
	w.RegisterTestCase(TC_manager_check_active_package_for_sim)

	/* Delete package for sim */
	w.RegisterTestCase(TC_manager_set_delete_package_for_sim)

	/* Get Packages after removal fo the a package*/
	w.RegisterTestCase(TC_manager_get_package_after_removal_for_sim)

	/* Inactivate sim */
	w.RegisterTestCase(TC_manager_inactivate_sim)

	/* Get Sim by subscriber after inactivation */
	w.RegisterTestCase(TC_manager_get_sim_by_subscriber)

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
