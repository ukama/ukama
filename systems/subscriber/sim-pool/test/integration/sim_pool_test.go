package integration

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	sb "github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/server"

	"google.golang.org/grpc"

	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSimServiceServer(s, &sb.SimPoolServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestAddSimToSimPool(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSimServiceClient(conn)
	_, err = client.Add(ctx, &pb.AddRequest{
		Sim: []*pb.AddSim{
			{
			Msisdn:"1234567890",
			IsPhysical:     true,
			ActivationCode: "123456",
			SmDpAddress:    "http://localhost:8080",
			QrCode:         "http://localhost:8080/qr/123456",
			SimType:        pb.SimType_INTER_MNO_DATA,
		},

			
		},
	})
	if err != nil {
		t.Fatalf("AddRequest failed: %v", err)
	}

	resp, err := client.Get(ctx, &pb.GetRequest{
		IsPhysicalSim: true,
		SimType:       pb.SimType_INTER_MNO_DATA,
	})
	if err != nil {
		t.Fatalf("GetRequest failed: %v", err)
	}
	expected := &pb.GetResponse{
		Sim: &pb.Sim{
			Id:             1,
			IsAllocated:    false,
			IsPhysical:     true,
			Msisdn:         "1234567890",
			ActivationCode: "123456",
			SmDpAddress:    "http://localhost:8080",
			QrCode:         "http://localhost:8080/qr/123456",
			SimType:        pb.SimType_INTER_MNO_DATA,
		},
	}
	if !cmp.Equal(resp, expected) {
		t.Errorf("Add Sim test failed, expected %v but got %v", expected, resp)
	}
}


func TestGetSimFromSimPool(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSimServiceClient(conn)

	resp, err := client.Get(ctx, &pb.GetRequest{
		IsPhysicalSim: true,
		SimType:       pb.SimType_INTER_MNO_DATA,
	})
	if err != nil {
		t.Fatalf("AddRequest failed: %v", err)
	}
	expected := &pb.GetResponse{
		Sim: &pb.Sim{
			Id:             1,
			IsAllocated:    false,
			IsPhysical:     true,
			Msisdn:         "1234567890",
			ActivationCode: "123456",
			Iccid:          "1234567890123456789",
			SmDpAddress:    "http://localhost:8080",
			QrCode:         "http://localhost:8080/qr/123456",
			SimType:        pb.SimType_INTER_MNO_DATA,
		},
	}
	if !cmp.Equal(resp, expected) {
		t.Errorf("Get Sim test failed, expected %v but got %v", expected, resp)
	}
}
