package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

type PackageServer struct {
	packageRepo db.PackageRepo
	pb.UnimplementedPackagesServiceServer
}

func NewPackageServer(packageRepo db.PackageRepo) *PackageServer {
	return &PackageServer{packageRepo: packageRepo}

}

func (p *PackageServer) GetPackage(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	logrus.Infof("Get rate %v", req.GetId())

	packageId := req.GetId()

	_package, err := p.packageRepo.GetPackage(packageId)
	if err != nil {
		logrus.Error("error while getting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}
	resp := &pb.GetPackageResponse{
		Package: dbPackageToPbPackages(_package),
	}

	return resp, nil
}

func (p *PackageServer) GetPackages(ctx context.Context, req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error) {
	logrus.Infof("GetPackages")
	packages, err := p.packageRepo.GetPackages()
	if err != nil {
		logrus.Error("error while getting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	packageList := &pb.GetPackagesResponse{
		Packages: dbpackagesToPbPackages(packages),
	}

	return packageList, nil
}

func (p *PackageServer) CreatePackage(ctx context.Context, req *pb.CreatePackageRequest) (*pb.CreatePackageResponse, error) {
	_package := db.Package{
		Name: req.GetName(),
		// Sim_type: validations.ReqPbToStr(pb.SimType(req.GetSimType())),
		// Sim_type:     validations.ReqPbToStr(req.GetSimType()),
		Org_id:       uint(req.GetOrgId()),
		Active:       req.Active,
		Duration:     uint(req.GetDuration()),
		Sms_volume:   uint(req.GetSmsVolume()),
		Data_volume:  uint(req.GetDataVolume()),
		Voice_volume: uint(req.GetVoiceVolume()),
		Org_rates_id: uint(req.GetOrgId()),
	}

	_packageRes, err := p.packageRepo.CreatePackage(_package)
	if err != nil {
		logrus.Error("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}

	return &pb.CreatePackageResponse{
		Package: dbPackageToPbPackages(&_packageRes),
	}, nil
}
func dbpackagesToPbPackages(packages []db.Package) []*pb.Package {
	res := []*pb.Package{}
	for _, u := range packages {
		res = append(res, dbPackageToPbPackages(&u))
	}
	return res
}

func dbPackageToPbPackages(r *db.Package) *pb.Package {
	return &pb.Package{
		Id:          uint64(r.ID),
		Name:        r.Name,
		OrgId:       int64(r.Org_id),
		Active:      r.Active,
		Duration:    uint64(r.Duration),
		SmsVolume:   int64(r.Sms_volume),
		OrgRatesId:  uint64(r.Org_rates_id),
		DataVolume:  int64(r.Data_volume),
		VoiceVolume: int64(r.Voice_volume),
		// SimType:     r.Sim_type.String(),
		CreatedAt: r.CreatedAt.String(),
		UpdatedAt: r.UpdatedAt.String(),
		DeletedAt: r.DeletedAt.Time.String(),
	}
}
