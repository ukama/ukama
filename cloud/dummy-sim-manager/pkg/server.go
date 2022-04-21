package pkg

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)

type SimManagerServer struct {
	simmgr.UnimplementedSimManagerServiceServer
	etcd *clientv3.Client
}

func NewSimManagerServer(etcdHost string) *SimManagerServer {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   []string{etcdHost},
	})
	if err != nil {
		logrus.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &SimManagerServer{
		etcd: client,
	}
}

const iccidPrefix = "890000"
const etcKeyPrifix = "dummy-sim-mgr."

type simInfo struct {
	Iccid    string           `json:"iccid"`
	Imsi     string           `json:"imsi"`
	Services *simmgr.Services `json:"services"`
}

func (s SimManagerServer) SetServiceStatus(ctx context.Context, request *simmgr.SetServiceStatusRequest) (*simmgr.SetServiceStatusResponse, error) {
	logrus.Infof("SetServiceStatus: %+v", request)
	if !strings.HasPrefix(request.Iccid, iccidPrefix) {
		return nil, status.Errorf(codes.NotFound, "Sim not found. Dummy sim should start with "+iccidPrefix)
	}

	sim := s.getSimInfo(ctx, request.Iccid)

	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "Sim not found.")
	}

	sim.Services = request.Services

	_, err := s.etcd.Put(ctx, getEtcdKey(request.Iccid), marshalSimInfo(sim))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update sim info in etcd: %v", err)
	}

	return &simmgr.SetServiceStatusResponse{}, nil
}

func (s SimManagerServer) getSimInfo(ctx context.Context, iccid string) *simInfo {
	val, err := s.etcd.Get(ctx, getEtcdKey(iccid), clientv3.WithLimit(1))
	if err != nil {
		logrus.Errorf("Cannot get sim info from etcd: %v", err)
		return nil
	}

	var sim *simInfo
	if val.Count > 0 {
		sim = unmarshalSimInfo(val.Kvs[0].Value)
	} else {
		sim = nil
	}
	return sim
}

func marshalSimInfo(info *simInfo) string {
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		logrus.Errorf("Cannot marshal sim info: %v", err)
	}
	return string(b)
}

func unmarshalSimInfo(b []byte) *simInfo {
	info := simInfo{}
	err := json.Unmarshal(b, &info)
	if err != nil {
		logrus.Errorf("Cannot unmarshal sim info: %v", err)
		return nil
	}
	return &info
}

func (s SimManagerServer) GetSimStatus(ctx context.Context, request *simmgr.GetSimStatusRequest) (*simmgr.GetSimStatusResponse, error) {
	logrus.Infof("GetSimStatus: %+v", request)
	if !strings.HasPrefix(request.Iccid, iccidPrefix) {
		return nil, status.Errorf(codes.NotFound, "Sim not found. Dummy sim should start with "+iccidPrefix)
	}

	sim := s.getSimInfo(ctx, request.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "Sim not found.")
	}

	return &simmgr.GetSimStatusResponse{
		Status:   simmgr.GetSimStatusResponse_ACTIVE,
		Services: sim.Services,
	}, nil
}

func (s SimManagerServer) GetSimInfo(ctx context.Context, request *simmgr.GetSimInfoRequest) (*simmgr.GetSimInfoResponse, error) {
	logrus.Infof("GetSimInfo: %+v", request)
	if !strings.HasPrefix(request.Iccid, iccidPrefix) {
		return nil, status.Errorf(codes.NotFound, "Sim not found. Dummy sim should start with "+iccidPrefix)
	}
	iccid := request.Iccid

	sim, err := s.getOrCreateSim(ctx, request, iccid)
	if err != nil {
		return nil, err
	}

	return &simmgr.GetSimInfoResponse{
		Iccid: sim.Iccid,
		Imsi:  sim.Imsi,
	}, nil
}

func (s SimManagerServer) getOrCreateSim(ctx context.Context, request *simmgr.GetSimInfoRequest, iccid string) (*simInfo, error) {
	logrus.Infof("Get sim info for iccid: %s", iccid)
	sim := s.getSimInfo(ctx, request.Iccid)
	if sim == nil {

		imsi := request.Iccid[len(iccid)-15:]
		sim = &simInfo{
			Iccid: iccid,
			Imsi:  imsi,
			Services: &simmgr.Services{
				Data:  &wrapperspb.BoolValue{Value: true},
				Sms:   &wrapperspb.BoolValue{Value: false},
				Voice: &wrapperspb.BoolValue{Value: false},
			},
		}
	}

	_, err := s.etcd.Put(ctx, getEtcdKey(request.Iccid), marshalSimInfo(sim))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update sim info in etcd: %v", err)
	}
	return sim, nil
}

func (s SimManagerServer) TerminateSim(ctx context.Context, request *simmgr.TerminateSimRequest) (*simmgr.TerminateSimResponse, error) {
	logrus.Infof("Terminate sim for iccid: %s", request.Iccid)
	sim := s.getSimInfo(ctx, request.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "Sim not found.")
	}

	_, err := s.etcd.Delete(ctx, getEtcdKey(request.Iccid))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot delete sim info from etcd: %v", err)
	}
	return &simmgr.TerminateSimResponse{}, nil
}

func getEtcdKey(key string) string {
	return etcKeyPrifix + key
}
