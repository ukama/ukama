package server

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb"
	validations "github.com/ukama/ukama/systems/data-plan/package/pkg/validations"

	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PackageServer struct {
	packageRepo db.PackageRepo
	pb.UnimplementedPackagesServiceServer
}

func NewPackageServer(packageRepo db.PackageRepo) *PackageServer {
	return &PackageServer{packageRepo: packageRepo}
}

func (p *PackageServer) Get(ctx context.Context, req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error) {
	logrus.Infof("GetPackages : %v  ,%v", req.GetOrgId(), req.GetId())

	if validations.IsInvalidId(req.GetOrgId()) {
		return nil, status.Errorf(codes.InvalidArgument, "OrgId is required.")
	}

	packages, err := p.packageRepo.Get(req.GetOrgId(), req.GetId())
	if err != nil {
		logrus.Error("error while getting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	packageList := &pb.GetPackagesResponse{
		Packages: dbpackagesToPbPackages(packages),
	}

	return packageList, nil
}

func (p *PackageServer) Add(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	logrus.Infof("Add Package Name: %v, SimType: %v, Active: %v, Duration: %v, SmsVolume: %v, DataVolume: %v, Voice_volume: %v", req.Name, req.SimType, req.Active, req.Duration, req.SmsVolume, req.DataVolume, req.VoiceVolume)
	_package := &db.Package{
		Name:         req.GetName(),
		Sim_type:     req.GetSimType().String(),
		Org_id:       uint(req.GetOrgId()),
		Active:       req.Active,
		Duration:     uint(req.GetDuration()),
		Sms_volume:   uint(req.GetSmsVolume()),
		Data_volume:  uint(req.GetDataVolume()),
		Voice_volume: uint(req.GetVoiceVolume()),
		Org_rates_id: uint(req.GetOrgRatesId()),
	}
	err := p.packageRepo.Add(_package)
	if err != nil {

		logrus.Error("Error adding a package. " + err.Error())

		return nil, status.Errorf(codes.Internal, "error adding a package")
	}

	return &pb.AddPackageResponse{Package: dbPackageToPbPackages(_package)}, nil

}

func (p *PackageServer) Delete(ctx context.Context, req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	logrus.Infof("Delete Packages, orgId: %v, packageId: %v", req.GetOrgId(), req.GetId())
	if validations.IsInvalidId(req.GetOrgId()) {
		return nil, status.Errorf(codes.InvalidArgument, "OrgId is required.")
	}
	if validations.IsInvalidId(req.GetId()) {
		return nil, status.Errorf(codes.InvalidArgument, "PackageID is required.")
	}
	err := p.packageRepo.Delete(req.GetOrgId(), req.GetId())
	if err != nil {
		logrus.Error("error while deleting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}
	return &pb.DeletePackageResponse{Id: strconv.FormatUint(req.Id, 10)}, nil
}

func (p *PackageServer) Update(ctx context.Context, req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	logrus.Infof("Update Package Id: %v, Name: %v, SimType: %v, Active: %v, Duration: %v, SmsVolume: %v, DataVolume: %v, Voice_volume: %v",
		req.Id, req.Name, req.SimType, req.Active, req.Duration, req.SmsVolume, req.DataVolume, req.VoiceVolume)
	_package := db.Package{
		Name:         req.GetName(),
		Sim_type:     req.GetSimType().String(),
		Active:       req.Active,
		Duration:     uint(req.GetDuration()),
		Sms_volume:   uint(req.GetSmsVolume()),
		Data_volume:  uint(req.GetDataVolume()),
		Voice_volume: uint(req.GetVoiceVolume()),
		Org_rates_id: uint(req.GetOrgRatesId()),
	}

	_packages, err := p.packageRepo.Update(req.Id, _package)
	if err != nil {
		logrus.Error("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
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
		Id:          uint64(p.ID),
		Name:        p.Name,
		OrgId:       int64(p.Org_id),
		Active:      p.Active,
		Duration:    uint64(p.Duration),
		SmsVolume:   int64(p.Sms_volume),
		OrgRatesId:  uint64(p.Org_rates_id),
		DataVolume:  int64(p.Data_volume),
		VoiceVolume: int64(p.Voice_volume),
		SimType:     pb.SimType(pb.SimType_value[p.Sim_type]),
		CreatedAt:   p.CreatedAt.String(),
		UpdatedAt:   p.UpdatedAt.String(),
		DeletedAt:   p.DeletedAt.Time.String(),
	}
}
