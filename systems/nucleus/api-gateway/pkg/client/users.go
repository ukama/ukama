/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pbusers "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
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
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	client := pbusers.NewUserServiceClient(conn)

	return &Users{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func (u *Users) Close() {
	if u.conn != nil {
		if err := u.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close User Service connection: %v", err)
		}
	}
}

func (u *Users) Get(userId string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return u.client.Get(ctx, &pbusers.GetRequest{UserId: userId})
}

func (u *Users) GetByEmail(email string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return u.client.GetByEmail(ctx, &pbusers.GetByEmailRequest{Email: email})
}

func (u *Users) GetByAuthId(authId string) (*pbusers.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return u.client.GetByAuthId(ctx, &pbusers.GetByAuthIdRequest{AuthId: authId})
}

func (u *Users) AddUser(user *pbusers.User) (*pbusers.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.client.Add(ctx, &pbusers.AddRequest{User: user})
}

func (u *Users) UpdateUser(userId string, user *pbusers.UserAttributes) (*pbusers.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return u.client.Update(ctx, &pbusers.UpdateRequest{
		UserId: userId,
		User: &pbusers.UserAttributes{
			Email: user.Email,
			Phone: user.Phone,
			Name:  user.Name,
		},
	})
}

func (u *Users) Delete(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	_, err := u.client.Delete(ctx, &pbusers.DeleteRequest{UserId: userId})
	return err
}

func (u *Users) DeactivateUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	_, err := u.client.Deactivate(ctx, &pbusers.DeactivateRequest{UserId: userId})
	return err
}

func (u *Users) Whoami(userId string) (*pbusers.WhoamiResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return u.client.Whoami(ctx, &pbusers.GetRequest{UserId: userId})
}
