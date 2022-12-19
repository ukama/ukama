package simPool

import (
	// pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"

	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"
)

type SimPoolServer struct {
	simPoolRepo db.SimPoolRepo
	// pb.UnimplementedSimPoolServiceServer
}

func NewSimPoolServer(simPoolRepo db.SimPoolRepo) *SimPoolServer {
	return &SimPoolServer{simPoolRepo: simPoolRepo}
}

// func (p *SimPoolServer) Get(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
// 	logrus.Infof("GetPackage : %v ", req.GetId())
// 	_package, err := p.simPoolRepo.Get(req.GetId())

// 	if err != nil {
// 		logrus.Error("error getting a package" + err.Error())

// 		return nil, grpc.SqlErrorToGrpc(err, "package")
// 	}

// 	resp := &pb.GetPackageResponse{SimPool: dbPackageToPbPackages(_package)}

// 	return resp, nil

// }

// func dbpackagesToPbPackages(packages []db.SimPool) []*pb.SimPool {
// 	res := []*pb.SimPool{}
// 	for _, u := range packages {
// 		res = append(res, dbPackageToPbPackages(&u))
// 	}
// 	return res
// }

// func dbPackageToPbPackages(p *db.SimPool) *pb.SimPool {
// 	return &pb.SimPool{}
// }
