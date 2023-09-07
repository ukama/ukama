package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Exporter struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.ExporterServiceClient
	host    string
}

func NewExporter(host string, timeout time.Duration) *Exporter {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewExporterServiceClient(conn)

	return &Exporter{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewExporterFromClient(c pb.ExporterServiceClient) *Exporter {
	return &Exporter{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (r *Exporter) Close() {
	r.conn.Close()
}

func (e *Exporter) Dummy(req *pb.DummyParameter) (*pb.DummyParameter, error) {
	_, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return &pb.DummyParameter{}, nil
}
