package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pbusers "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Users struct {
	conn    *grpc.ClientConn
	client  pbusers.UserServiceClient
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
	client := pbusers.NewUserServiceClient(conn)

	return &Users{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewTestHssFromClient(networkClient pbusers.UserServiceClient) *Users {
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

func (r *Users) AddUser(user *pbusers.User, requesterId string) (*pbusers.AddResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Add(ctx, &pbusers.AddRequest{User: user})
}

func (r *Users) UpdateUser(userUUID string, user *pbusers.UserAttributes, requesterId string) (*pbusers.UpdateResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Update(ctx, &pbusers.UpdateRequest{
		UserUuid: userUUID,
		User: &pbusers.UserAttributes{
			Email: user.Email,
			Phone: user.Phone,
			Name:  user.Name,
		},
	})
}

func (r *Users) Delete(userUUID string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.Delete(ctx, &pbusers.DeleteRequest{UserUuid: userUUID})
	return err
}

func (r *Users) Get(userUUID string, requesterId string) (*pbusers.GetResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Get(ctx, &pbusers.GetRequest{UserUuid: userUUID})
}

func (r *Users) DeactivateUser(userUUID string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.Deactivate(ctx, &pbusers.DeactivateRequest{UserUuid: userUUID})
	return err
}

func (r *Users) getContext(requester string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	md := metadata.Pairs("x-requester", requester)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
