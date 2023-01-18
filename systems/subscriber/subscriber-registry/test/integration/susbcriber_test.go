package integration

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	sb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/server"

	"google.golang.org/grpc"

	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSubscriberRegistryServiceServer(s, &sb.SubcriberServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAddSubscriber(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSubscriberRegistryServiceClient(conn)
	resp, err := client.Add(ctx, &pb.AddSubscriberRequest{
		FirstName:             "John",
		LastName:              "Doe",
		Email:                 "johndoe@example.com",
		PhoneNumber:           "+1234567890",
		Address:               "123 Main St.",
		IdSerial:              "123456789",
		NetworkID:             "00000000-0000-0000-0000-000000000000",
		ProofOfIdentification: "Drivers License 12345678",
	})
	if err != nil {
		t.Fatalf("AddRequest failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	expected := &pb.AddSubscriberResponse{
		Subscriber: &pb.Subscriber{
			FirstName:             "John",
			LastName:              "Doe",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "+1234567890",
			Address:               "123 Main St.",
			IdSerial:              "123456789",
			NetworkID:             "00000000-0000-0000-0000-000000000000",
			ProofOfIdentification: "Drivers License 12345678",
		},
	}
	if !cmp.Equal(resp, expected) {
		t.Errorf("Add subscriber test failed, expected %v but got %v", expected, resp)
	}
}
func TestUpdateSubscriber(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSubscriberRegistryServiceClient(conn)

	// Add a subscriber before updating
	_, err = client.Add(ctx, &pb.AddSubscriberRequest{
		FirstName:             "John",
		LastName:              "Doe",
		Email:                 "johndoe@example.com",
		PhoneNumber:           "+1234567890",
		Address:               "123 Main St.",
		IdSerial:              "123456789",
		NetworkID:             "00000000-0000-0000-0000-000000000000",
		ProofOfIdentification: "Drivers License 12345678",
	})
	if err != nil {
		t.Fatalf("AddRequest failed: %v", err)
	}

	// Perform update request
	resp, err := client.Update(ctx, &pb.UpdateSubscriberRequest{
		IdSerial:              "123456789",
		Email:                 "janedoe@example.com",
		PhoneNumber:           "+0987654321",
		Address:               "456 Park Ave.",
		ProofOfIdentification: "Passport 987654321",
	})
	if err != nil {
		t.Fatalf("UpdateRequest failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	expected := &pb.UpdateSubscriberResponse{
		IdSerial:              "123456789",
		Email:                 "janedoe@example.com",
		PhoneNumber:           "+0987654321",
		Address:               "456 Park Ave.",
		ProofOfIdentification: "Passport 987654321",
	}
	if !cmp.Equal(resp, expected) {
		t.Errorf("Update subscriber test failed, expected %v but got %v", expected, resp)
	}

}
func TestDeleteSubscriber(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSubscriberRegistryServiceClient(conn)

	// Add a subscriber before deleting
	_, err = client.Add(ctx, &pb.AddSubscriberRequest{
		FirstName:             "John",
		LastName:              "Doe",
		Email:                 "johndoe@example.com",
		PhoneNumber:           "+1234567890",
		Address:               "123 Main St.",
		IdSerial:              "123456789",
		NetworkID:             "00000000-0000-0000-0000-000000000000",
		ProofOfIdentification: "Drivers License 12345678",
	})
	if err != nil {
		t.Fatalf("AddRequest failed: %v", err)
	}

	_, err = client.Delete(ctx, &pb.DeleteSubscriberRequest{
		SubscriberID: "123456789",
	})
	if err != nil {
		t.Fatalf("DeleteRequest failed: %v", err)
	}

}
