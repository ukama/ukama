package server

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/ukama/ukama/systems/common/sql"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber/pkg/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
	var gormdb sql.Db;
	srv := NewSubscriberServer(db.NewSubscriberRepo(gormdb))

	pb.RegisterSubscriberServiceServer(s, srv)
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
	
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestSayHello(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    client := pb.NewSubscriberServiceClient(conn)
	req:=pb.AddSubscriberRequest{
		FirstName: "John",
		LastName: "Doe",
		Email: "john.doe@example.com",
		PhoneNumber: "123-456-7890",
		Address: "123 Main St.",
		Gender: "male"}

    resp, err := client.Add(ctx, &req)
    if err != nil {
        t.Fatalf("SayHello failed: %v", err)
    }
    log.Printf("Response: %+v", resp)
    // Test for output here.
}