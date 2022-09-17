package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Lookup struct {
	conn    *grpc.ClientConn
	client  pb.UserServiceClient
	timeout time.Duration
	host    string
}

func Newlookup(host string, timeout time.Duration) *Lookup {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewLookupServiceClient(conn)

	return &Lookup{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func (r *Lookup) Close() {
	r.conn.Close()
}

func (r *Lookup) getContext(requester string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	md := metadata.Pairs("x-requester", requester)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}

func (r *Lookup) AddOrg(orgName string, user *pb.User, simToken string, requesterId string) (*pb.AddOrgResponse, error) {
	, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.AddOrg(ctx, &pb.AddOrgRequest{Org: orgName, User: user, SimToken: simToken})
}
