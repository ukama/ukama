package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/orchestrator/constructor/pb/gen"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/contractor"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConstructorServer struct {
	debug          bool
	oRepo          db.OrgRepo
	dRepo          db.DeploymentRepo
	cRepo          db.ConfigRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedConstructorServiceServer
}

func NewConstructorServer(o db.OrgRepo, d db.DeploymentRepo, c db.ConfigRepo, msgBus mb.MsgBusServiceClient, debug bool) *ConstructorServer {
	return &ConstructorServer{
		debug:          debug,
		dRepo:          d,
		oRepo:          o,
		cRepo:          c,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (c *ConstructorServer) ConstructOrg(ctx context.Context, in *pb.ConstructOrgRequest) (*pb.ConstructOrgResponse, error) {
	return &pb.ConstructOrgResponse{}, nil
}

func (c *ConstructorServer) DistructOrg(ctx context.Context, in *pb.DistructOrgRequest) (*pb.DistructOrgResponse, error) {
	return &pb.DistructOrgResponse{}, nil
}

func (c *ConstructorServer) Deployment(ctx context.Context, in *pb.DeploymentRequest) (*pb.DeploymentResponse, error) {
	return &pb.DeploymentResponse{}, nil
}

func (c *ConstructorServer) GetDeployment(ctx context.Context, in *pb.GetDeploymentRequests) (*pb.GetDeploymentResponse, error) {
	return &pb.GetDeploymentResponse{}, nil
}

func (c *ConstructorServer) RemoveDeployment(ctx context.Context, in *pb.RemoveDeploymentRequest) (*pb.RemoveDeploymentResponse, error) {
	return &pb.RemoveDeploymentResponse{}, nil
}

func (c *ConstructorServer) AddConfig(ctx context.Context, in *pb.AddConfigRequest) (*pb.AddConfigResponse, error) {

	if len(in.Config.KeyVal)%2 != 0 {
		log.Errorf("Invalid config for key value pair : %+v", in.Config.KeyVal)
		return nil, status.Error(codes.InvalidArgument, "invalid key value pair")
	}

	cfg := pbConfigTodb(in.Config)

	/* Get Config first */
	oCfg, err := c.cRepo.Get(in.Config.Name)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			log.Errorf("Failed to get system config: %v", err)
			return nil, grpc.SqlErrorToGrpc(err, "config")
		}
	}

	/* Create a new/updated config */
	err = c.cRepo.Create(cfg)
	if err != nil {
		log.Errorf("Failed to create system config: %v", err)
		return nil, grpc.SqlErrorToGrpc(err, "config")
	}

	/* Start the procedure to update older versions of the orgs */
	l, err := c.deployConfigToOrgs(*cfg, oCfg.Orgs)
	if err != nil {
		log.Errorf("Failed to delpoy config: %v", err)
		return nil, err
	}

	return &pb.AddConfigResponse{
		OrgId: l,
	}, nil
}

func (c *ConstructorServer) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	cfg, err := c.cRepo.Get(in.GetName())
	if err != nil {
		log.Errorf("Failed to get %s system config: %v", in.GetName(), err)
		return nil, grpc.SqlErrorToGrpc(err, "config")
	}

	return &pb.GetConfigResponse{
		Config: dbConfigToPb(cfg),
	}, nil
}

func (c *ConstructorServer) GetDeploymentHistory(ctx context.Context, in *pb.GetDeploymentRequests) (*pb.GetDeploymentResponse, error) {
	return &pb.GetDeploymentResponse{}, nil
}

func pbConfigTodb(c *pb.Config) *db.Config {
	return &db.Config{
		Name:    c.Name,
		Source:  c.Source,
		Chart:   c.Chart,
		Version: c.Version,
		Values:  c.KeyVal,
	}
}

func dbConfigToPb(c *db.Config) *pb.Config {
	return &pb.Config{
		Name:    c.Name,
		Source:  c.Source,
		Chart:   c.Chart,
		Version: c.Version,
		KeyVal:  c.Values,
	}
}

func (c *ConstructorServer) deploy(ctx context.Context, cfg db.Config, org db.Org) {

	kv := CreateKeyValuePair(cfg, org)

	name := org.OrgName + "-" + cfg.Name
	d := &db.Deployment{
		Name:      name,
		Namespace: name,
		Code:      uuid.NewV4(),
		Org:       org,
		Status:    db.Waiting,
		Values:    MapToArray(kv),
		Config:    cfg,
	}

	ct, err := contractor.NewContractor(kv, c.debug, cfg.Source, d.Namespace)
	if err != nil {
		log.Errorf("Failed to start the deployment %s for %s org and system %s. Error %v", d.Name, org.OrgId, cfg.Name, err)
		return
	}

	err = ct.ApplyHelmfile(ctx)
	if err != nil {
		log.Errorf("Failed to deploy %s for %s org and system %s. Error %v", d.Name, org.OrgId, cfg.Name, err)
		return
	}

	err = c.dRepo.Add(d)
	if err != nil {
		log.Errorf("Failed to update deployment %s for %s org and system %s in Db. Error %v", d.Name, org.OrgId, cfg.Name, err)
		return
	}
}

func (c *ConstructorServer) deployConfigToOrgs(cfg db.Config, orgs []db.Org) ([]string, error) {
	resp := make([]string, 0, len(orgs))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	for _, org := range orgs {
		go c.deploy(ctx, cfg, org)
		resp = append(resp, org.OrgId.String())
	}
	return resp, nil
}

func CreateKeyValuePair(cfg db.Config, org db.Org) map[string]interface{} {

	kv1 := ArrayToMap(cfg.Values)

	kv2 := ArrayToMap(org.Values)

	kv := MergeMaps(kv1, kv2)

	return kv
}

func ArrayToMap(arr []string) map[string]string {
	kv := make(map[string]string)
	for i := 0; i < len(arr); i += 2 {
		kv[arr[i]] = arr[i+1]
	}
	return kv
}

/* If similar keys values of m2 will override m1 */
func MergeMaps(m1 map[string]string, m2 map[string]string) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

func MapToArray(m map[string]interface{}) []string {
	arr := make([]string, 0, len(m)*2)
	i := 0
	for k, v := range m {
		arr[i] = k
		arr[i+1] = v.(string)
		i += 2
	}
	return arr
}
