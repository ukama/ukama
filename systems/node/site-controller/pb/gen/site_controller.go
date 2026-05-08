package gen

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SiteControllerServiceName = "ukama.node.site_controller.v1.SiteControllerService"

type SetSiteRequest struct {
	SiteId string `json:"site_id"`
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}
type SetSiteResponse struct {
	State *SiteState `json:"state,omitempty"`
}
type SetServiceRequest struct {
	SiteId string `json:"site_id"`
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}
type SetServiceResponse struct {
	State *SiteState `json:"state,omitempty"`
}
type SetRadioRequest struct {
	SiteId string `json:"site_id"`
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}
type SetRadioResponse struct {
	State *SiteState `json:"state,omitempty"`
}
type GetSiteStateRequest struct {
	SiteId string `json:"site_id"`
}
type GetSiteStateResponse struct {
	State *SiteState `json:"state,omitempty"`
}
type UpsertPortMapRequest struct {
	SiteId  string          `json:"site_id"`
	CnodeId string          `json:"cnode_id,omitempty"`
	Ports   []*PortMapEntry `json:"ports"`
}
type UpsertPortMapResponse struct{}
type GetPortMapRequest struct {
	SiteId string `json:"site_id"`
}
type GetPortMapResponse struct {
	Ports []*PortMapEntry `json:"ports"`
}
type ApplySwitchPolicyRequest struct {
	SiteId string `json:"site_id"`
}
type ApplySwitchPolicyResponse struct {
	Applied bool `json:"applied"`
}
type PowerCycleNodeRequest struct {
	SiteId string `json:"site_id"`
	Role   string `json:"role"`
	Reason string `json:"reason,omitempty"`
}
type PowerCycleNodeResponse struct{}
type PortMapEntry struct {
	Port   int32  `json:"port"`
	Role   string `json:"role"`
	NodeId string `json:"node_id,omitempty"`
	Class  string `json:"class"`
	Policy string `json:"policy"`
}
type SiteState struct {
	SiteId         string `json:"site_id"`
	DesiredSite    string `json:"desired_site"`
	DesiredService string `json:"desired_service"`
	DesiredRadio   string `json:"desired_radio"`
	Power          string `json:"power"`
	Service        string `json:"service"`
	Radio          string `json:"radio"`
	Access         string `json:"access"`
	Reason         string `json:"reason"`
}

type SiteControllerServiceClient interface {
	SetSite(context.Context, *SetSiteRequest, ...grpc.CallOption) (*SetSiteResponse, error)
	SetService(context.Context, *SetServiceRequest, ...grpc.CallOption) (*SetServiceResponse, error)
	SetRadio(context.Context, *SetRadioRequest, ...grpc.CallOption) (*SetRadioResponse, error)
	GetSiteState(context.Context, *GetSiteStateRequest, ...grpc.CallOption) (*GetSiteStateResponse, error)
	UpsertPortMap(context.Context, *UpsertPortMapRequest, ...grpc.CallOption) (*UpsertPortMapResponse, error)
	GetPortMap(context.Context, *GetPortMapRequest, ...grpc.CallOption) (*GetPortMapResponse, error)
	ApplySwitchPolicy(context.Context, *ApplySwitchPolicyRequest, ...grpc.CallOption) (*ApplySwitchPolicyResponse, error)
	PowerCycleNode(context.Context, *PowerCycleNodeRequest, ...grpc.CallOption) (*PowerCycleNodeResponse, error)
}
type siteControllerServiceClient struct{ cc grpc.ClientConnInterface }

func NewSiteControllerServiceClient(cc grpc.ClientConnInterface) SiteControllerServiceClient {
	return &siteControllerServiceClient{cc}
}
func (c *siteControllerServiceClient) invoke(ctx context.Context, method string, in, out interface{}, opts ...grpc.CallOption) error {
	opts = append(opts, grpc.ForceCodec(ForceJSONCodec()))
	return c.cc.Invoke(ctx, "/"+SiteControllerServiceName+"/"+method, in, out, opts...)
}
func (c *siteControllerServiceClient) SetSite(ctx context.Context, in *SetSiteRequest, opts ...grpc.CallOption) (*SetSiteResponse, error) {
	out := new(SetSiteResponse)
	err := c.invoke(ctx, "SetSite", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) SetService(ctx context.Context, in *SetServiceRequest, opts ...grpc.CallOption) (*SetServiceResponse, error) {
	out := new(SetServiceResponse)
	err := c.invoke(ctx, "SetService", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) SetRadio(ctx context.Context, in *SetRadioRequest, opts ...grpc.CallOption) (*SetRadioResponse, error) {
	out := new(SetRadioResponse)
	err := c.invoke(ctx, "SetRadio", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) GetSiteState(ctx context.Context, in *GetSiteStateRequest, opts ...grpc.CallOption) (*GetSiteStateResponse, error) {
	out := new(GetSiteStateResponse)
	err := c.invoke(ctx, "GetSiteState", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) UpsertPortMap(ctx context.Context, in *UpsertPortMapRequest, opts ...grpc.CallOption) (*UpsertPortMapResponse, error) {
	out := new(UpsertPortMapResponse)
	err := c.invoke(ctx, "UpsertPortMap", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) GetPortMap(ctx context.Context, in *GetPortMapRequest, opts ...grpc.CallOption) (*GetPortMapResponse, error) {
	out := new(GetPortMapResponse)
	err := c.invoke(ctx, "GetPortMap", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) ApplySwitchPolicy(ctx context.Context, in *ApplySwitchPolicyRequest, opts ...grpc.CallOption) (*ApplySwitchPolicyResponse, error) {
	out := new(ApplySwitchPolicyResponse)
	err := c.invoke(ctx, "ApplySwitchPolicy", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) PowerCycleNode(ctx context.Context, in *PowerCycleNodeRequest, opts ...grpc.CallOption) (*PowerCycleNodeResponse, error) {
	out := new(PowerCycleNodeResponse)
	err := c.invoke(ctx, "PowerCycleNode", in, out, opts...)
	return out, err
}

type SiteControllerServiceServer interface {
	SetSite(context.Context, *SetSiteRequest) (*SetSiteResponse, error)
	SetService(context.Context, *SetServiceRequest) (*SetServiceResponse, error)
	SetRadio(context.Context, *SetRadioRequest) (*SetRadioResponse, error)
	GetSiteState(context.Context, *GetSiteStateRequest) (*GetSiteStateResponse, error)
	UpsertPortMap(context.Context, *UpsertPortMapRequest) (*UpsertPortMapResponse, error)
	GetPortMap(context.Context, *GetPortMapRequest) (*GetPortMapResponse, error)
	ApplySwitchPolicy(context.Context, *ApplySwitchPolicyRequest) (*ApplySwitchPolicyResponse, error)
	PowerCycleNode(context.Context, *PowerCycleNodeRequest) (*PowerCycleNodeResponse, error)
}
type UnimplementedSiteControllerServiceServer struct{}

func (UnimplementedSiteControllerServiceServer) SetSite(context.Context, *SetSiteRequest) (*SetSiteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetSite not implemented")
}
func (UnimplementedSiteControllerServiceServer) SetService(context.Context, *SetServiceRequest) (*SetServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetService not implemented")
}
func (UnimplementedSiteControllerServiceServer) SetRadio(context.Context, *SetRadioRequest) (*SetRadioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRadio not implemented")
}
func (UnimplementedSiteControllerServiceServer) GetSiteState(context.Context, *GetSiteStateRequest) (*GetSiteStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSiteState not implemented")
}
func (UnimplementedSiteControllerServiceServer) UpsertPortMap(context.Context, *UpsertPortMapRequest) (*UpsertPortMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertPortMap not implemented")
}
func (UnimplementedSiteControllerServiceServer) GetPortMap(context.Context, *GetPortMapRequest) (*GetPortMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortMap not implemented")
}
func (UnimplementedSiteControllerServiceServer) ApplySwitchPolicy(context.Context, *ApplySwitchPolicyRequest) (*ApplySwitchPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplySwitchPolicy not implemented")
}
func (UnimplementedSiteControllerServiceServer) PowerCycleNode(context.Context, *PowerCycleNodeRequest) (*PowerCycleNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PowerCycleNode not implemented")
}
func RegisterSiteControllerServiceServer(s grpc.ServiceRegistrar, srv SiteControllerServiceServer) {
	s.RegisterService(&SiteControllerService_ServiceDesc, srv)
}
func unary[T any](srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor, method string, call func(SiteControllerServiceServer, context.Context, *T) (interface{}, error)) (interface{}, error) {
	in := new(T)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(srv.(SiteControllerServiceServer), ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + SiteControllerServiceName + "/" + method}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return call(srv.(SiteControllerServiceServer), ctx, req.(*T))
	}
	return interceptor(ctx, in, info, handler)
}

var SiteControllerService_ServiceDesc = grpc.ServiceDesc{ServiceName: SiteControllerServiceName, HandlerType: (*SiteControllerServiceServer)(nil), Methods: []grpc.MethodDesc{
	{MethodName: "SetSite", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[SetSiteRequest](s, c, d, i, "SetSite", func(srv SiteControllerServiceServer, ctx context.Context, req *SetSiteRequest) (interface{}, error) {
			return srv.SetSite(ctx, req)
		})
	}},
	{MethodName: "SetService", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[SetServiceRequest](s, c, d, i, "SetService", func(srv SiteControllerServiceServer, ctx context.Context, req *SetServiceRequest) (interface{}, error) {
			return srv.SetService(ctx, req)
		})
	}},
	{MethodName: "SetRadio", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[SetRadioRequest](s, c, d, i, "SetRadio", func(srv SiteControllerServiceServer, ctx context.Context, req *SetRadioRequest) (interface{}, error) {
			return srv.SetRadio(ctx, req)
		})
	}},
	{MethodName: "GetSiteState", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[GetSiteStateRequest](s, c, d, i, "GetSiteState", func(srv SiteControllerServiceServer, ctx context.Context, req *GetSiteStateRequest) (interface{}, error) {
			return srv.GetSiteState(ctx, req)
		})
	}},
	{MethodName: "UpsertPortMap", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[UpsertPortMapRequest](s, c, d, i, "UpsertPortMap", func(srv SiteControllerServiceServer, ctx context.Context, req *UpsertPortMapRequest) (interface{}, error) {
			return srv.UpsertPortMap(ctx, req)
		})
	}},
	{MethodName: "GetPortMap", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[GetPortMapRequest](s, c, d, i, "GetPortMap", func(srv SiteControllerServiceServer, ctx context.Context, req *GetPortMapRequest) (interface{}, error) {
			return srv.GetPortMap(ctx, req)
		})
	}},
	{MethodName: "ApplySwitchPolicy", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[ApplySwitchPolicyRequest](s, c, d, i, "ApplySwitchPolicy", func(srv SiteControllerServiceServer, ctx context.Context, req *ApplySwitchPolicyRequest) (interface{}, error) {
			return srv.ApplySwitchPolicy(ctx, req)
		})
	}},
	{MethodName: "PowerCycleNode", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return unary[PowerCycleNodeRequest](s, c, d, i, "PowerCycleNode", func(srv SiteControllerServiceServer, ctx context.Context, req *PowerCycleNodeRequest) (interface{}, error) {
			return srv.PowerCycleNode(ctx, req)
		})
	}}}, Streams: []grpc.StreamDesc{}, Metadata: "site_controller.json.grpc"}
