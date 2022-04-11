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

func (r *Hss) AddUser(orgName string, user *pb.User, simToken string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	return r.client.Add(ctx, &pb.AddRequest{Org: orgName, User: user, SimToken: simToken})
}

func (r *Hss) UpdateUser(userId string, user *pb.UserAttributes) (*pb.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	return r.client.Update(ctx, &pb.UpdateRequest{
		Uuid: userId,
		User: &pb.UserAttributes{
			Email: user.Email,
			Phone: user.Phone,
			Name:  user.Name,
		},
	})
}

func (r *Hss) GetUsers(orgName string) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	return r.client.List(ctx, &pb.ListRequest{
		Org: orgName,
	})
}

func (r *Hss) Delete(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	_, err := r.client.Delete(ctx, &pb.DeleteRequest{Uuid: userId})
	return err
}

func (r *Hss) Get(userId string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	defer cancel()

	return r.client.Get(ctx, &pb.GetRequest{Uuid: userId})
}
