package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Users struct {
	conn    *grpc.ClientConn
	client  pb.UserServiceClient
	timeout time.Duration
	host    string
}

func NewUsers(host string, timeout time.Duration) *Users {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewUserServiceClient(conn)

	return &Users{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewTestHssFromClient(networkClient pb.UserServiceClient) *Users {
	return &Users{
		host:    "localhost",
		timeout: 1,
		conn:    nil,
		client:  networkClient,
	}
}

func (r *Users) Close() {
	r.conn.Close()
}

func (r *Users) AddUser(orgName string, user *pb.User, simToken string, requesterId string) (*pb.AddResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Add(ctx, &pb.AddRequest{Org: orgName, User: user, SimToken: simToken})
}

func (r *Users) UpdateUser(userId string, user *pb.UserAttributes, requesterId string) (*pb.UpdateResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Update(ctx, &pb.UpdateRequest{
		UserId: userId,
		User: &pb.UserAttributes{
			Email: user.Email,
			Phone: user.Phone,
			Name:  user.Name,
		},
	})
}

func (r *Users) GetUsers(orgName string, requesterId string) (*pb.ListResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.List(ctx, &pb.ListRequest{
		Org: orgName,
	})
}

func (r *Users) Delete(userId string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.Delete(ctx, &pb.DeleteRequest{UserId: userId})
	return err
}

func (r *Users) Get(userId string, requesterId string) (*pb.GetResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Get(ctx, &pb.GetRequest{UserId: userId})
}

func (r *Users) SetSimStatus(req *pb.SetSimStatusRequest, requesterId string) (*pb.Sim, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	resp, err := r.client.SetSimStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Sim, err
}

func (r *Users) DeactivateUser(userId string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.DeactivateUser(ctx, &pb.DeactivateUserRequest{UserId: userId})
	return err
}

func (r *Users) GetQr(iccid string, requesterId string) (*pb.GetQrCodeResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.GetQrCode(ctx, &pb.GetQrCodeRequest{Iccid: iccid})
}

func (r *Users) getContext(requester string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	md := metadata.Pairs("x-requester", requester)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
