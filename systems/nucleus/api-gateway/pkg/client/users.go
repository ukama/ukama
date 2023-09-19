package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	pbusers "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Users struct {
	conn    *grpc.ClientConn
	client  pbusers.UserServiceClient
	timeout time.Duration
	host    string
}

func NewUserRegistryFromClient(client pbusers.UserServiceClient) *Users {
	return &Users{
		timeout: 1 * time.Second,
		conn:    nil,
		client:  client,
	}
}

func NewUsers(host string, timeout time.Duration) *Users {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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

func (r *Users) Get(userId string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Get(ctx, &pbusers.GetRequest{UserId: userId})
}

func (r *Users) GetByEmail(email string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetByEmail(ctx, &pbusers.GetByEmailRequest{Email: email})
}

func (r *Users) GetByAuthId(authId string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetByAuthId(ctx, &pbusers.GetByAuthIdRequest{AuthId: authId})
}

func (r *Users) AddUser(user *pbusers.User) (*pbusers.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.Add(ctx, &pbusers.AddRequest{User: user})
}

func (r *Users) UpdateUser(userId string, user *pbusers.UserAttributes) (*pbusers.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
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

func (r *Users) Delete(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.Delete(ctx, &pbusers.DeleteRequest{UserId: userId})
	return err
}

func (r *Users) DeactivateUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.Deactivate(ctx, &pbusers.DeactivateRequest{UserId: userId})
	return err
}

func (r *Users) Whoami(userId string) (*pbusers.WhoamiResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Whoami(ctx, &pbusers.GetRequest{UserId: userId})
}
