package server

import (
	"context"
	"testing"

	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"google.golang.org/grpc"
)

func TestSimPoolServer_Add(t *testing.T) {
    // Create a gRPC client to connect to the server
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()
    client := pb.NewSimServiceClient(conn)
    // Create a request for a physical SIM of type "SIM_TYPE_1"
    // req := &pb.GetRequest{IsPhysicalSim: true, SimType: pb.SimType_SIM_TYPE_1}
	req:=&pb.AddRequest{
	Sim:[]*pb.AddSim{
		{
			Iccid: "123456789", SimType: pb.SimType_INTER_MNO_DATA, Msisdn: "555-555-1234", SmDpAddress: "http://example.com", ActivationCode: "abc123", QrCode: "qr123", IsPhysical: true,
		},
		{
			Iccid: "12273", SimType: pb.SimType_INTER_MNO_DATA, Msisdn: "583-5343-0234", SmDpAddress: "http://example.com", ActivationCode: "abc123", QrCode: "qr123", IsPhysical: true,
		},
	},
}


    // Send the request to the server and check the response
    resp, err := client.Add(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to adding SIMs: %v", err)
    }
    if resp.Sim == nil {
        t.Error("Expected a SIM but got nil")
    }
}


func TestSimPoolServer_Get(t *testing.T) {
    // Create a gRPC client to connect to the server
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()
    client := pb.NewSimServiceClient(conn)
	req:=&pb.GetRequest{
		IsPhysicalSim:true,
		SimType:pb.SimType_INTER_MNO_DATA,
}


    // Send the request to the server and check the response
    resp, err := client.Get(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to get SIMs: %v", err)
    }
    if resp.Sim == nil {
        t.Error("Expected a SIM but got nil")
    }
}
func TestSimPoolServer_GetByICCID(t *testing.T) {
    // Create a gRPC client to connect to the server
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()
    client := pb.NewSimServiceClient(conn)
	req:=&pb.GetByIccidRequest{
		Iccid:"123456789",
}


    // Send the request to the server and check the response
    resp, err := client.GetByIccid(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to get SIM by ICCID: %v", err)
    }
    if resp.Sim == nil {
        t.Error("Expected a SIM but got nil")
    }
}

func TestSimPoolServer_Delete(t *testing.T) {
    // Create a gRPC client to connect to the server
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()
    client := pb.NewSimServiceClient(conn)
	req:=&pb.DeleteRequest{
		Id:[]uint64{
			123456789,
			123456789,
		},
}

    // Send the request to the server and check the response
    resp, err := client.Delete(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to delete SIMs with: %v", err)
    }
    if resp.Id == nil {
        t.Error("Expected Id but got nil")
    }
}
