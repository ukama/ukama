/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"time"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/validation"
	"github.com/ukama/ukama/systems/data-plan/package/pkg"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/client"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
)

type PackageServer struct {
	orgName        string
	packageRepo    db.PackageRepo
	rate           client.RateClientProvider
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedPackagesServiceServer
	orgId string
}

func NewPackageServer(orgName string, packageRepo db.PackageRepo, rate client.RateClientProvider, msgBus mb.MsgBusServiceClient, orgId string) *PackageServer {
	return &PackageServer{
		orgName:        orgName,
		packageRepo:    packageRepo,
		msgbus:         msgBus,
		rate:           rate,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		orgId:          orgId,
	}

}

func (p *PackageServer) Get(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	log.Infof("GetPackage : %v ", packageID)
	_package, err := p.packageRepo.Get(packageID)

	if err != nil {
		log.Error("error getting a package" + err.Error())

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

	log.Infof("GetPackageDetails : %v ", packageID)
	_package, err := p.packageRepo.GetDetails(packageID)

	if err != nil {
		log.Error("error getting a package" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	resp := &pb.GetPackageResponse{Package: dbPackageToPbPackages(_package)}

	return resp, nil
}

func (p *PackageServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	log.Infof("Get all packages")

	packages, err := p.packageRepo.GetAll()
	if err != nil {
		log.Error("error while getting package by Org" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "packages")
	}

	packageList := &pb.GetAllResponse{
		Packages: dbpackagesToPbPackages(packages),
	}

	return packageList, nil
}

func (p *PackageServer) Add(ctx context.Context, req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {

	log.Infof("Adding package %v", req)

	ownId, err := uuid.FromString(req.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of owner uuid. Error %s", err.Error())
	}

	baserate, err := uuid.FromString(req.GetBaserateId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of base rate. Error %s", err.Error())
	}

	formattedFrom, err := validation.ValidateDate(req.From)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	}
	if err := validation.IsFutureDate(formattedFrom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())

	}

	formattedTo, err := validation.ValidateDate(req.To)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	}
	if err := validation.IsFutureDate(formattedTo); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())

	}
	if err := validation.IsAfterDate(formattedTo, formattedFrom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())

	}

	from, err := validation.FromString(formattedFrom)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	}

	to, err := validation.FromString(formattedTo)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error: %s", err.Error())
	}

	pkgUuid := uuid.NewV4()

	pkg := db.Package{
		Uuid:         pkgUuid,
		OwnerId:      ownId,
		Name:         req.GetName(),
		SimType:      ukama.ParseSimType(req.GetSimType()),
		Active:       req.Active,
		From:         from,
		To:           to,
		Duration:     uint64(req.GetDuration()),
		SmsVolume:    uint64(req.GetSmsVolume()),
		DataVolume:   uint64(req.GetDataVolume()),
		VoiceVolume:  uint64(req.GetVoiceVolume()),
		MessageUnits: ukama.ParseMessageType(req.MessageUnit),
		VoiceUnits:   ukama.ParseCallUnitType(req.VoiceUnit),
		DataUnits:    ukama.ParseDataUnitType(req.DataUnit),
		Flatrate:     req.Flatrate,
		Type:         ukama.ParsePackageType(req.Type),
		PackageMarkup: db.PackageMarkup{
			BaseRateId: baserate,
			Markup:     req.Markup,
		},
		PackageDetails: db.PackageDetails{
			Apn: req.Apn,
		},
		Currency:      req.Currency,
		Overdraft:     req.Overdraft,
		TrafficPolicy: req.TrafficPolicy,
		Networks:      req.Networks,
		Country:       req.Country,
		SyncStatus:    ukama.StatusTypePending,
	}

	pr := db.PackageRate{
		Amount:    req.Amount,
		SmsMo:     0.0,
		SmsMt:     0.0,
		Data:      0.0,
		PackageID: pkgUuid,
	}

	rateSvc, err := p.rate.GetClient()
	if err != nil {
		return nil, err
	}

	rate, err := rateSvc.GetRateById(ctx, &rpb.GetRateByIdRequest{
		OwnerId:  req.OwnerId,
		BaseRate: req.BaserateId,
	})
	if err != nil {
		log.Errorf("Failed to get base rate for package. Error: %s", err.Error())
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid base id. Error %s", err.Error())
	}

	pkg.Country = rate.Rate.Country
	pkg.Provider = rate.Rate.Provider

	if pkg.PackageDetails.Apn == "" {
		pkg.PackageDetails.Apn = rate.Rate.Apn
	}

	/* Only when package is not fixed anount */
	if !pkg.Flatrate {
		// calculae rate per unit
		calculateRatePerUnit(&pkg.PackageRate, rate.Rate, pkg.MessageUnits, pkg.DataUnits)
		calculateTotalAmount(&pkg)
	}

	err = p.packageRepo.Add(&pkg, &pr)
	if err != nil {
		log.Error("Error while adding a package. " + err.Error())
		return nil, status.Errorf(codes.Internal, "Error while adding a package.")
	}

	resp := &pb.AddPackageResponse{Package: dbPackageToPbPackages(&pkg)}

	if p.msgbus != nil {
		route := p.baseRoutingKey.SetActionCreate().SetObject("package").MustBuild()
		evt := &epb.CreatePackageEvent{
			Uuid:            resp.Package.Uuid,
			OrgId:           p.orgId,
			OwnerId:         resp.Package.OwnerId,
			Type:            resp.Package.Type,
			Flatrate:        resp.Package.Flatrate,
			Amount:          resp.Package.Rate.Amount,
			From:            resp.Package.From,
			To:              resp.Package.To,
			SimType:         resp.Package.SimType,
			SmsVolume:       resp.Package.SmsVolume,
			DataVolume:      resp.Package.DataVolume,
			VoiceVolume:     resp.Package.VoiceVolume,
			DataUnit:        resp.Package.DataUnit,
			VoiceUnit:       resp.Package.VoiceUnit,
			Messageunit:     resp.Package.MessageUnit,
			DataUnitCost:    pkg.PackageRate.Data,
			MessageUnitCost: pkg.PackageRate.SmsMo,
			VoiceUnitCost:   pkg.PackageRate.SmsMt,
		}
		err = p.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}

	}
	return resp, nil
}

func (p *PackageServer) Delete(ctx context.Context, req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	log.Infof("Delete Packages packageId: %v", req.GetUuid())

	err = p.packageRepo.Delete(packageID)
	if err != nil {
		log.Error("error while deleting package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if p.msgbus != nil {
		evt := &epb.DeletePackageEvent{
			Uuid:  req.Uuid,
			OrgId: p.orgId,
		}
		route := p.baseRoutingKey.SetActionDelete().SetObject("package").MustBuild()
		err = p.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.DeletePackageResponse{
		Uuid: req.GetUuid(),
	}, nil
}

func (p *PackageServer) Update(ctx context.Context, req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	log.Infof("Update Package Uuid: %v, Name: %v,Active: %v",
		req.Uuid, req.Name, req.Active)
	_package := &db.Package{
		Name:   req.GetName(),
		Active: req.Active,
	}

	packageID, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of package uuid. Error %s", err.Error())
	}

	err = p.packageRepo.Update(packageID, _package)
	if err != nil {
		log.Error("error while getting updating a package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	updatedPackage, err := p.packageRepo.Get(packageID)
	if err != nil {
		log.Error("error while getting updated package" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	if p.msgbus != nil {
		route := p.baseRoutingKey.SetAction("update").SetObject("package").MustBuild()
		evt := &epb.UpdatePackageEvent{
			Uuid:  req.Uuid,
			OrgId: p.orgId,
		}
		err = p.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.UpdatePackageResponse{Package: dbPackageToPbPackages(updatedPackage)}, nil
}

func dbpackagesToPbPackages(packages []db.Package) []*pb.Package {
	res := []*pb.Package{}
	for _, u := range packages {
		res = append(res, dbPackageToPbPackages(&u))
	}
	return res
}

func dbPackageToPbPackages(p *db.Package) *pb.Package {
	var d string
	if p.DeletedAt.Valid {
		d = p.DeletedAt.Time.Format(time.RFC3339)
	}

	return &pb.Package{
		Uuid:        p.Uuid.String(),
		Name:        p.Name,
		Active:      p.Active,
		From:        p.From.Format(time.RFC3339),
		To:          p.To.Format(time.RFC3339),
		SmsVolume:   int64(p.SmsVolume),
		DataVolume:  int64(p.DataVolume),
		VoiceVolume: int64(p.VoiceVolume),
		SimType:     p.SimType.String(),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   d,
		Rate: &pb.PackageRate{
			Data:   p.PackageRate.Data,
			SmsMo:  p.PackageRate.SmsMo,
			SmsMt:  p.PackageRate.SmsMt,
			Amount: p.PackageRate.Amount,
		},
		Markup: &pb.PackageMarkup{
			Baserate: p.PackageMarkup.BaseRateId.String(),
			Markup:   p.PackageMarkup.Markup,
		},
		OwnerId:       p.OwnerId.String(),
		Provider:      p.Provider,
		Duration:      p.Duration,
		Amount:        p.PackageRate.Amount,
		Type:          p.Type.String(),
		MessageUnit:   p.MessageUnits.String(),
		VoiceUnit:     p.VoiceUnits.String(),
		DataUnit:      p.DataUnits.String(),
		Country:       p.Country,
		Currency:      p.Currency,
		Apn:           p.PackageDetails.Apn,
		Flatrate:      p.Flatrate,
		Overdraft:     p.Overdraft,
		TrafficPolicy: p.TrafficPolicy,
		Networks:      p.Networks,
		SyncStatus:    p.SyncStatus.String(),
	}
}

func calculateRatePerUnit(pr *db.PackageRate, rate *bpb.Rate, mu ukama.MessageUnitType, du ukama.DataUnitType) {

	pr.SmsMo = (float64)(ukama.ReturnMessageUnits(mu)) * rate.SmsMo
	pr.SmsMt = (float64)(ukama.ReturnMessageUnits(mu)) * rate.SmsMt
	pr.Data = (float64)(ukama.ReturnDataUnits(du)) * rate.Data

}

func calculateTotalAmount(pr *db.Package) {

	pr.PackageRate.Amount = (pr.PackageRate.SmsMo * float64(pr.SmsVolume)) +
		(pr.PackageRate.SmsMt * float64(pr.SmsVolume)) +
		(pr.PackageRate.Data * float64(pr.DataVolume))

}
