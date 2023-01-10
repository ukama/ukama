package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSubcriberServer_Add(t *testing.T) {

	testCases := []struct {
		name         string
		req          *pb.AddSubscriberRequest
		expectedResp *pb.AddSubscriberResponse
		expectedErr  error
	}{
		{
			name: "Valid request",
			req: &pb.AddSubscriberRequest{
				FirstName:             "John",
				LastName:              "Doe",
				NetworkID:             "00000000-0000-0000-0000-000000000000",
				Email:                 "john.doe@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "male",
				Address:               "123 Main St",
				IdSerial:              "123456",
				ProofOfIdentification: "drivers_license",
			},
			expectedResp: &pb.AddSubscriberResponse{
				Subscriber: &pb.Subscriber{
					FirstName:             "John",
					LastName:              "Doe",
					NetworkID:             "00000000-0000-0000-0000-000000000000",
					Email:                 "john.doe@example.com",
					PhoneNumber:           "1234567890",
					Gender:                "male",
					Address:               "123 Main St",
					IdSerial:              "123456",
					ProofOfIdentification: "drivers_license",
				},
			},

			expectedErr: nil,
		},
		{
			name: "Invalid network ID",
			req: &pb.AddSubscriberRequest{
				FirstName:             "John",
				LastName:              "Doe",
				NetworkID:             "invalid",
				Email:                 "john.doe@example.com",
				PhoneNumber:           "1234567890",
				Gender:                "male",
				Address:               "123 Main St",
				DateOfBirth:           &timestamppb.Timestamp{Seconds: time.Now().Unix()},
				IdSerial:              "123456",
				ProofOfIdentification: "drivers_license",
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.InvalidArgument, "invalid UUID length: invalid"),
		},
		{
			name: "Error from repository",
			req: &pb.AddSubscriberRequest{
				FirstName: "John",
				LastName:  "Derick",
			},
			expectedResp: nil,
			expectedErr:  status.Error(codes.Internal, "repository error"),
		},
	}

	// Run test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			subscriberRepo := &subscriberRepoMock{
				addFunc: func(subscriber *db.Subscriber) error {
					if subscriber.NetworkID.String() == "00000000-0000-0000-0000-000000000000" {
						subscriberID := uuid.FromStringOrNil("93a3f36b-4556-444f-97c3-6132e0bfdda9")
						subscriber.SubscriberID = subscriberID
						return nil
					}
					return status.Error(codes.Internal, "repository error")
				},
				delFunc: func(subscriberID uuid.UUID) error {
					return nil
				},
				getFunc: func(subscriberID uuid.UUID) (*db.Subscriber, error) {
					return nil, nil
				},
				getByNetwork: func(networkID uuid.UUID) ([]db.Subscriber, error) {
					return nil, nil
				},
			}
			server := NewSubscriberServer(subscriberRepo)
			resp, err := server.Add(context.Background(), testCase.req)

			assert.Equal(t, testCase.expectedResp, resp)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

type subscriberRepoMock struct {
	addFunc      func(subscriber *db.Subscriber) error
	delFunc      func(subscriberID uuid.UUID) error
	getFunc      func(subscriberID uuid.UUID) (*db.Subscriber, error)
	getByNetwork func(networkID uuid.UUID) ([]db.Subscriber, error)
	updateFunc   func(subscriberID uuid.UUID, sub db.Subscriber) (db.Subscriber, error)
}

func (r *subscriberRepoMock) Update(subscriberID uuid.UUID, sub db.Subscriber) (*db.Subscriber, error) {
	if r.updateFunc != nil {
		updatedSubscriber, err := r.updateFunc(subscriberID, sub)
		return &updatedSubscriber, err
	}
	return nil, fmt.Errorf("updateFunc is not defined")
}

func (r *subscriberRepoMock) Add(subscriber *db.Subscriber) error {
	if r.addFunc != nil {
		return r.addFunc(subscriber)
	}
	return nil
}

func (r *subscriberRepoMock) Delete(subscriberID uuid.UUID) error {
	if r.delFunc != nil {
		return r.delFunc(subscriberID)
	}
	return nil
}

func (r *subscriberRepoMock) Get(subscriberID uuid.UUID) (*db.Subscriber, error) {
	if r.getFunc != nil {
		return r.getFunc(subscriberID)
	}
	return nil, nil
}
func (r *subscriberRepoMock) GetByNetwork(networkID uuid.UUID) ([]db.Subscriber, error) {
	if r.getByNetwork != nil {
		return r.getByNetwork(networkID)
	}
	return nil, nil
}
