package client

import (
	"context"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Hss struct {
	conn    *grpc.ClientConn
	client  pb.UserServiceClient
	timeout int
	host    string
}

func NewHss(host string, timeout int) *Hss {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewUserServiceClient(conn)

	return &Hss{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewTestHssFromClient(registryClient pb.UserServiceClient) *Hss {
	return &Hss{
		host:    "localhost",
		timeout: 1,
		conn:    nil,
		client:  registryClient,
	}
}

func (r *Hss) Close() {
	r.conn.Close()
}

func (r *Hss) AddUser(orgName string, user *pb.User) (*pb.AddUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.Add(ctx, &pb.AddUserRequest{Org: orgName, User: user})
	if grpcErr, ok := marshalError(err); ok {
		return nil, grpcErr
	}
	return res, nil
}

func (r *Hss) GetUsers(orgName string) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	res, err := r.client.List(ctx, &pb.ListUsersRequest{
		Org: orgName,
	})
	if grpcErr, ok := marshalError(err); ok {
		return nil, grpcErr
	}
	return res, nil
}

func (r *Hss) Delete(orgName string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	_, err := r.client.Delete(ctx, &pb.DeleteUserRequest{UserUuid: userId, Org: orgName})
	return err
}
