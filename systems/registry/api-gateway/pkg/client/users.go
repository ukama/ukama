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

func (r *Users) Close() {
	r.conn.Close()
}

func (r *Users) Get(userId string, requesterId string) (*pbusers.GetResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Get(ctx, &pbusers.GetRequest{UserId: userId})
}

func (r *Users) GetByAuthId(authId string, requesterId string) (*pbusers.GetResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.GetByAuthId(ctx, &pbusers.GetByAuthIdRequest{AuthId: authId})
}

func (r *Users) AddUser(user *pbusers.User, requesterId string) (*pbusers.AddResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Add(ctx, &pbusers.AddRequest{User: user})
}

func (r *Users) UpdateUser(userId string, user *pbusers.UserAttributes, requesterId string) (*pbusers.UpdateResponse, error) {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	return r.client.Update(ctx, &pbusers.UpdateRequest{
		UserId: userId,
		User: &pbusers.UserAttributes{
			Email: user.Email,
			Phone: user.Phone,
			Name:  user.Name,
		},
	})
}

func (r *Users) Delete(userId string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.Delete(ctx, &pbusers.DeleteRequest{UserId: userId})
	return err
}

func (r *Users) DeactivateUser(userId string, requesterId string) error {
	ctx, cancel := r.getContext(requesterId)
	defer cancel()

	_, err := r.client.Deactivate(ctx, &pbusers.DeactivateRequest{UserId: userId})
	return err
}

func (r *Users) getContext(requester string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.timeout)*time.Second)
	md := metadata.Pairs("x-requester", requester)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
