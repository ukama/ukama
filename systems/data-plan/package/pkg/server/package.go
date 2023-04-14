package server

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/client"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PackageServer struct {
	packageRepo    db.PackageRepo
	rate           client.RateService
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedPackagesServiceServer
}

func NewPackageServer(packageRepo db.PackageRepo, rate client.RateService, msgBus mb.MsgBusServiceClient) *PackageServer {
	return &PackageServer{
		packageRepo:    packageRepo,
		msgbus:         msgBus,
		rate:           rate,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (p *PackageServer) Get(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	logrus.Infof("GetPackage : %v ", packageID)
	_package, err := p.packageRepo.Get(packageID)

	if err != nil {
		logrus.Error("error getting a package" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	resp := &pb.GetPackageResponse{Package: dbPackageToPbPackages(_package)}

	return resp, nil
}

func (p *PackageServer) GetDetails(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	logrus.Infof("GetPackage : %v ", packageID)
	_package, err := p.packageRepo.GetDetails(packageID)

	if err != nil {
		logrus.Error("error getting a package" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	resp := &pb.GetPackageResponse{Package: dbPackageToPbPackages(_package)}

	return resp, nil
}

func (p *PackageServer) GetByOrg(ctx context.Context, req *pb.GetByOrgPackageRequest) (*pb.GetByOrgPackageResponse, error) {
	logrus.Infof("GetPackage by Org: %v ", req.GetOrgId())

	orgId, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	packages, err := p.packageRepo.GetByOrg(orgId)
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

	// Need to get Network for ukama_data
	// Possible activation date
	// What happens if we have two rates for a period

	orgId, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	_package := &db.Package{
		Uuid:         uuid.NewV4(),
		Name:         req.GetName(),
		SimType:      db.ParseType(req.GetSimType()),
		OrgId:        orgId,
		Active:       req.Active,
		Duration:     uint64(req.GetDuration()),
		SmsVolume:    uint64(req.GetSmsVolume()),
		DataVolume:   uint64(req.GetDataVolume()),
		VoiceVolume:  uint64(req.GetVoiceVolume()),
		MessageUnits: db.ParseMessageType(req.Messageunit),
		CallUnits:    db.ParseCallUnitType(req.CallUnit),
		DataUnits:    db.ParseDataUnitType(req.DataUnit),
		Dlbr:         uint64(req.Dlbr),
		Ulbr:         uint64(req.Ulbr),
		Flatrate:     req.Flatrate,
		Rate: &db.PackageRate{
			Flatrate: req.FlatrateAmount,
		},
		Markup: &db.PackageMarkup{
			Markup: req.Markup,
		},
	}

	// Request rate
	rate, err := p.rate.GetRate(req.Baserate)
	if err != nil {
		log.Errorf("Failed to get base rate for package. Error: %s", err.Error()))
		retuun nil, tatus.Errorf(codes.InvalidArgument,
			"invalid base id. Error %s", err.Error())
	}

	// calculae rate per unit

	err = p.packageRepo.Add(_package)
	if err != nil {

		logrus.Error("Error while adding a package. " + err.Error())
		return nil, status.Errorf(codes.Internal, "Error while adding a package.")
	}
	return &pb.AddPackageResponse{Package: dbPackageToPbPackages(_package)}, nil
}

func (p *PackageServer) Delete(ctx context.Context, req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	logrus.Infof("Delete Packages packageId: %v", req.GetUuid())

	err = p.packageRepo.Delete(packageID)
	if err != nil {
		logrus.Error("error while deleting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	route := p.baseRoutingKey.SetAction("delete").SetObject("package").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.DeletePackageResponse{
		Uuid: req.GetUuid(),
	}, nil
}

func (p *PackageServer) Update(ctx context.Context, req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	logrus.Infof("Update Package Uuid: %v, Name: %v, SimType: %v, Active: %v, Duration: %v, SmsVolume: %v, DataVolume: %v, Voice_volume: %v",
		req.Uuid, req.Name, req.SimType, req.Active, req.Duration, req.SmsVolume, req.DataVolume, req.VoiceVolume)
	_package := &db.Package{
		Name:         req.GetName(),
		SimType:      db.ParseType(req.GetSimType()),
		Active:       req.Active,
		Duration:     uint64(req.GetDuration()),
		SmsVolume:    uint64(req.GetSmsVolume()),
		DataVolume:   uint64(req.GetDataVolume()),
		VoiceVolume:  uint64(req.GetVoiceVolume()),
		MessageUnits: db.ParseMessageType(req.Messageunit),
		CallUnits:    db.ParseCallUnitType(req.CallUnit),
		DataUnits:    db.ParseDataUnitType(req.DataUnit),
	}

	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	err = p.packageRepo.Update(packageID, _package)
	if err != nil {
		logrus.Error("error while getting updating a package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	route := p.baseRoutingKey.SetAction("update").SetObject("package").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.UpdatePackageResponse{Package: dbPackageToPbPackages(_package)}, nil
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
		Uuid:        p.Uuid.String(),
		Name:        p.Name,
		OrgId:       p.OrgId.String(),
		Active:      p.Active,
		Duration:    uint64(p.Duration),
		SmsVolume:   int64(p.SmsVolume),
		DataVolume:  int64(p.DataVolume),
		VoiceVolume: int64(p.VoiceVolume),
		SimType:     p.SimType.String(),
		CreatedAt:   p.CreatedAt.String(),
		UpdatedAt:   p.UpdatedAt.String(),
		DeletedAt:   p.DeletedAt.Time.String(),
		X2G:         p.X2g,
		X3G:         p.X3g,
		X5G:         p.X5g,
		Lte:         p.Lte,
		LteM:        p.LteM,
		Apn:         p.Apn,
		Vpmn:        p.Vpmn,
		Imsi:        p.Imsi,
		Rate: &pb.PackageRates{
			Data:     p.Rate.Data,
			SmsMo:    p.Rate.SmsMo,
			SmsMt:    p.Rate.SmsMt,
			FlatRate: p.Rate.Flatrate,
		},
		Markup: &pb.PackageMarkup{
			Baserate: p.Markup.BaseRateId.String(),
			Markup:   p.Markup.Markup,
		},
		Provider:    p.Provider,
		Dlbr:        int64(p.Dlbr),
		Ulbr:        int64(p.Ulbr),
		Messageunit: p.MessageUnits.String(),
		CallUnit:    p.CallUnits.String(),
		DataUnit:    p.DataUnits.String(),
		Country:     p.Country,
		Currency:    p.Currency,
		Flatrate:    p.Flatrate,
		EffectiveAt: p.EffectiveAt.Format(time.RFC3339),
		EndAt:       p.EndAt.Format(time.RFC3339),
	}
}
