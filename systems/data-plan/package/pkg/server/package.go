package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mbc "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type PackageServer struct {
	packageRepo   db.PackageRepo
	msgbus               mbc.MsgBusServiceClient
	packageRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedPackagesServiceServer
}

func NewPackageServer(packageRepo db.PackageRepo, msgBus mbc.MsgBusServiceClient) *PackageServer {
	return &PackageServer{
		packageRepo: packageRepo,
		msgbus:               msgBus,
		packageRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (p *PackageServer) Get(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	logrus.Infof("GetPackage : %v ", req.GetPackageID())
	packageID, err := uuid.FromString(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}
	_package, err := p.packageRepo.Get(packageID)

	if err != nil {
		logrus.Error("error getting a package" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	resp := &pb.GetPackageResponse{Package: dbPackageToPbPackages(_package)}

	return resp, nil
}
func (p *PackageServer) GetByOrg(ctx context.Context, req *pb.GetByOrgPackageRequest) (*pb.GetByOrgPackageResponse, error) {
	logrus.Infof("GetPackage by Org: %v ", req.GetOrgID())
	orgID, err := uuid.FromString(req.GetOrgID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	packages, err := p.packageRepo.GetByOrg(orgID)
	if err != nil {
		logrus.Error("error while getting package by Org" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	packageList := &pb.GetByOrgPackageResponse{
		Packages: dbpackagesToPbPackages(packages),
	}

	return packageList, nil
}
func (p *PackageServer) Add(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	logrus.Infof("Add Package Name: %v, SimType: %v, Active: %v, Duration: %v, SmsVolume: %v, DataVolume: %v, Voice_volume: %v", req.Name, req.SimType, req.Active, req.Duration, req.SmsVolume, req.DataVolume, req.VoiceVolume)
	orgID, err := uuid.FromString(req.GetOrgID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"error while parsing org as uuid. Error %s", err.Error())
	}
	strType := strings.ToLower(fmt.Sprintf("%v", req.GetSimType()))
	simType := db.ParseType(strType)
	
	if uint8(simType) != uint8(req.SimType) {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: provided sim type (%s) does not match with package allowed sim type (%s)",
			simType.String(), req.SimType)
	}
	
	_package := &db.Package{
		PackageID:         uuid.NewV4(),
		Name:         req.GetName(),
		SimType:      db.ParseType(req.GetSimType().String()),
		OrgID:       orgID,
		Active:       req.Active,
		Duration:     uint(req.GetDuration()),
		SmsVolume:   uint(req.GetSmsVolume()),
		DataVolume:  uint(req.GetDataVolume()),
		VoiceVolume: uint(req.GetVoiceVolume()),
		OrgRatesID: uint(req.GetOrgRatesID()),
	}
	err = p.packageRepo.Add(_package)
	if err != nil {

		logrus.Error("Error while adding a package. " + err.Error())

		return nil, status.Errorf(codes.Internal, "error adding a package")
	}

	return &pb.AddPackageResponse{Package: dbPackageToPbPackages(_package)}, nil
}

func (p *PackageServer) Delete(ctx context.Context, req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	logrus.Infof("Delete Packages packageId: %v", req.GetPackageID())
	packageID, err := uuid.FromString(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}
	err = p.packageRepo.Delete(packageID)
	if err != nil {
		logrus.Error("error while deleting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	// Publish message to msgbus

	route := p.packageRoutingKey.SetActionUpdate().SetObject("package").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.DeletePackageResponse{
	}, nil
}

func (p *PackageServer) Update(ctx context.Context, req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	logrus.Infof("Update Package Uuid: %v, Name: %v, SimType: %v, Active: %v, Duration: %v, SmsVolume: %v, DataVolume: %v, Voice_volume: %v",
		req.PackageID, req.Name, req.SimType, req.Active, req.Duration, req.SmsVolume, req.DataVolume, req.VoiceVolume)
		strType := strings.ToLower(fmt.Sprintf("%v", req.GetSimType()))
	simType := db.ParseType(strType)
	
	if uint8(simType) != uint8(req.SimType) {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: provided sim type (%s) does not match with package allowed sim type (%s)",
			simType.String(), req.SimType)
	}
	_package := db.Package{
		Name:         req.GetName(),
		SimType:     db.ParseType(req.GetSimType().String()),
		Active:       req.Active,
		Duration:     uint(req.GetDuration()),
		SmsVolume:   uint(req.GetSmsVolume()),
		DataVolume:  uint(req.GetDataVolume()),
		VoiceVolume: uint(req.GetVoiceVolume()),
		OrgRatesID: uint(req.GetOrgRatesID()),
	}
	packageID, err := uuid.FromString(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}
	_packages, err := p.packageRepo.Update(packageID, _package)
	if err != nil {
		logrus.Error("error while getting updating a package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	// Publish message to msgbus

	route := p.packageRoutingKey.SetActionUpdate().SetObject("package").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.UpdatePackageResponse{
		Package: dbPackageToPbPackages(_packages),
	}, nil
}

func dbpackagesToPbPackages(packages []db.Package) []*pb.Package {
	res := []*pb.Package{}
	for _, u := range packages {
		res = append(res, dbPackageToPbPackages(&u))
	}
	return res
}

func dbPackageToPbPackages(p *db.Package) *pb.Package {
	return &pb.Package{
		PackageID:p.PackageID.String(),
		Name:        p.Name,
		OrgID:       p.OrgID.String(),
		Active:      p.Active,
		Duration:    uint64(p.Duration),
		SmsVolume:   int64(p.SmsVolume),
		OrgRatesID:  uint64(p.OrgRatesID),
		DataVolume:  int64(p.DataVolume),
		VoiceVolume: int64(p.VoiceVolume),
		SimType:    pb.SimType(p.SimType),
		CreatedAt:   p.CreatedAt.String(),
		UpdatedAt:   p.UpdatedAt.String(),
		DeletedAt:   p.DeletedAt.Time.String(),
	}
}
