package rpc

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const CommandServiceName = "ukama.node.controller.v1.CommandService"

type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error)   { return json.Marshal(v) }
func (jsonCodec) Unmarshal(b []byte, v interface{}) error { return json.Unmarshal(b, v) }
func (jsonCodec) Name() string                            { return "json" }
func init()                                               { encoding.RegisterCodec(jsonCodec{}) }
func ForceJSONCodec() jsonCodec                           { return jsonCodec{} }

type SendNodeCommandRequest struct {
	NodeId string `json:"node_id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   []byte `json:"body,omitempty"`
	Source string `json:"source,omitempty"`
	Reason string `json:"reason,omitempty"`
}
type SendNodeCommandResponse struct{}
type CommandServiceClient interface {
	SendNodeCommand(context.Context, *SendNodeCommandRequest, ...grpc.CallOption) (*SendNodeCommandResponse, error)
}
type commandServiceClient struct{ cc grpc.ClientConnInterface }

func NewCommandServiceClient(cc grpc.ClientConnInterface) CommandServiceClient {
	return &commandServiceClient{cc: cc}
}
func (c *commandServiceClient) SendNodeCommand(ctx context.Context, in *SendNodeCommandRequest, opts ...grpc.CallOption) (*SendNodeCommandResponse, error) {
	out := new(SendNodeCommandResponse)
	opts = append(opts, grpc.ForceCodec(ForceJSONCodec()))
	err := c.cc.Invoke(ctx, "/"+CommandServiceName+"/SendNodeCommand", in, out, opts...)
	return out, err
}

type CommandServiceServer interface {
	SendNodeCommand(context.Context, *SendNodeCommandRequest) (*SendNodeCommandResponse, error)
}
type UnimplementedCommandServiceServer struct{}

func (UnimplementedCommandServiceServer) SendNodeCommand(context.Context, *SendNodeCommandRequest) (*SendNodeCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNodeCommand not implemented")
}
func RegisterCommandServiceServer(s grpc.ServiceRegistrar, srv CommandServiceServer) {
	s.RegisterService(&CommandService_ServiceDesc, srv)
}
func handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendNodeCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServiceServer).SendNodeCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + CommandServiceName + "/SendNodeCommand"}
	h := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServiceServer).SendNodeCommand(ctx, req.(*SendNodeCommandRequest))
	}
	return interceptor(ctx, in, info, h)
}

var CommandService_ServiceDesc = grpc.ServiceDesc{ServiceName: CommandServiceName, HandlerType: (*CommandServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "SendNodeCommand", Handler: handler}}, Streams: []grpc.StreamDesc{}, Metadata: "command.json.grpc"}
