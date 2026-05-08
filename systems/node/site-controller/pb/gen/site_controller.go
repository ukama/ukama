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
type GetSwitchPolicyRequest struct {
	SiteId string `json:"site_id"`
}
type GetSwitchPolicyResponse struct {
	Policy *SwitchPolicyStatus `json:"policy,omitempty"`
}
type RefreshSwitchPolicyRequest struct {
	SiteId  string `json:"site_id"`
	CnodeId string `json:"cnode_id,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
type RefreshSwitchPolicyResponse struct {
	Requested bool `json:"requested"`
}
type ReportSwitchPolicyRequest struct {
	SiteId  string        `json:"site_id"`
	CnodeId string        `json:"cnode_id"`
	Policy  *SwitchPolicy `json:"policy,omitempty"`
}
type ReportSwitchPolicyResponse struct {
	Policy *SwitchPolicyStatus `json:"policy,omitempty"`
}
type PowerCycleNodeRequest struct {
	SiteId string `json:"site_id"`
	Role   string `json:"role"`
	Reason string `json:"reason,omitempty"`
}
type PowerCycleNodeResponse struct{}

type SiteState struct {
	SiteId         string              `json:"site_id"`
	DesiredSite    string              `json:"desired_site"`
	DesiredService string              `json:"desired_service"`
	DesiredRadio   string              `json:"desired_radio"`
	Power          string              `json:"power"`
	Service        string              `json:"service"`
	Radio          string              `json:"radio"`
	Access         string              `json:"access"`
	Reason         string              `json:"reason"`
	SwitchPolicy   *SwitchPolicyStatus `json:"switch_policy,omitempty"`
}

type SwitchPolicyStatus struct {
	SiteId  string        `json:"site_id"`
	CnodeId string        `json:"cnode_id"`
	State   string        `json:"state"`
	Hash    string        `json:"hash"`
	Source  string        `json:"source"`
	Error   string        `json:"error"`
	Valid   bool          `json:"valid"`
	Reason  string        `json:"reason"`
	Policy  *SwitchPolicy `json:"policy,omitempty"`
}

type SwitchPolicy struct {
	SiteId    string              `json:"site_id"`
	Source    string              `json:"source"`
	UpdatedAt string              `json:"updated_at"`
	State     string              `json:"state,omitempty"`
	Hash      string              `json:"hash,omitempty"`
	Error     string              `json:"error,omitempty"`
	Ports     []*SwitchPolicyPort `json:"ports"`
}

type SwitchPolicyPort struct {
	Port   int32  `json:"port"`
	Role   string `json:"role"`
	NodeId string `json:"node_id,omitempty"`
	Class  string `json:"class"`
	Policy string `json:"policy"`
}

type SiteControllerServiceClient interface {
	SetSite(context.Context, *SetSiteRequest, ...grpc.CallOption) (*SetSiteResponse, error)
	SetService(context.Context, *SetServiceRequest, ...grpc.CallOption) (*SetServiceResponse, error)
	SetRadio(context.Context, *SetRadioRequest, ...grpc.CallOption) (*SetRadioResponse, error)
	GetSiteState(context.Context, *GetSiteStateRequest, ...grpc.CallOption) (*GetSiteStateResponse, error)
	GetSwitchPolicy(context.Context, *GetSwitchPolicyRequest, ...grpc.CallOption) (*GetSwitchPolicyResponse, error)
	RefreshSwitchPolicy(context.Context, *RefreshSwitchPolicyRequest, ...grpc.CallOption) (*RefreshSwitchPolicyResponse, error)
	ReportSwitchPolicy(context.Context, *ReportSwitchPolicyRequest, ...grpc.CallOption) (*ReportSwitchPolicyResponse, error)
	PowerCycleNode(context.Context, *PowerCycleNodeRequest, ...grpc.CallOption) (*PowerCycleNodeResponse, error)
}

type siteControllerServiceClient struct{ cc grpc.ClientConnInterface }

func NewSiteControllerServiceClient(cc grpc.ClientConnInterface) SiteControllerServiceClient {
	return &siteControllerServiceClient{cc: cc}
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
func (c *siteControllerServiceClient) GetSwitchPolicy(ctx context.Context, in *GetSwitchPolicyRequest, opts ...grpc.CallOption) (*GetSwitchPolicyResponse, error) {
	out := new(GetSwitchPolicyResponse)
	err := c.invoke(ctx, "GetSwitchPolicy", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) RefreshSwitchPolicy(ctx context.Context, in *RefreshSwitchPolicyRequest, opts ...grpc.CallOption) (*RefreshSwitchPolicyResponse, error) {
	out := new(RefreshSwitchPolicyResponse)
	err := c.invoke(ctx, "RefreshSwitchPolicy", in, out, opts...)
	return out, err
}
func (c *siteControllerServiceClient) ReportSwitchPolicy(ctx context.Context, in *ReportSwitchPolicyRequest, opts ...grpc.CallOption) (*ReportSwitchPolicyResponse, error) {
	out := new(ReportSwitchPolicyResponse)
	err := c.invoke(ctx, "ReportSwitchPolicy", in, out, opts...)
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
	GetSwitchPolicy(context.Context, *GetSwitchPolicyRequest) (*GetSwitchPolicyResponse, error)
	RefreshSwitchPolicy(context.Context, *RefreshSwitchPolicyRequest) (*RefreshSwitchPolicyResponse, error)
	ReportSwitchPolicy(context.Context, *ReportSwitchPolicyRequest) (*ReportSwitchPolicyResponse, error)
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
func (UnimplementedSiteControllerServiceServer) GetSwitchPolicy(context.Context, *GetSwitchPolicyRequest) (*GetSwitchPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSwitchPolicy not implemented")
}
func (UnimplementedSiteControllerServiceServer) RefreshSwitchPolicy(context.Context, *RefreshSwitchPolicyRequest) (*RefreshSwitchPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshSwitchPolicy not implemented")
}
func (UnimplementedSiteControllerServiceServer) ReportSwitchPolicy(context.Context, *ReportSwitchPolicyRequest) (*ReportSwitchPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportSwitchPolicy not implemented")
}
func (UnimplementedSiteControllerServiceServer) PowerCycleNode(context.Context, *PowerCycleNodeRequest) (*PowerCycleNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PowerCycleNode not implemented")
}

func RegisterSiteControllerServiceServer(s grpc.ServiceRegistrar, srv SiteControllerServiceServer) {
	s.RegisterService(&SiteControllerService_ServiceDesc, srv)
}

type scHandlerFunc func(SiteControllerServiceServer, context.Context, interface{}) (interface{}, error)

func handler[T any](srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor, method string, call func(SiteControllerServiceServer, context.Context, *T) (interface{}, error)) (interface{}, error) {
	in := new(T)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(srv.(SiteControllerServiceServer), ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + SiteControllerServiceName + "/" + method}
	h := func(ctx context.Context, req interface{}) (interface{}, error) {
		return call(srv.(SiteControllerServiceServer), ctx, req.(*T))
	}
	return interceptor(ctx, in, info, h)
}

var SiteControllerService_ServiceDesc = grpc.ServiceDesc{ServiceName: SiteControllerServiceName, HandlerType: (*SiteControllerServiceServer)(nil), Methods: []grpc.MethodDesc{
	{MethodName: "SetSite", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[SetSiteRequest](s, c, d, i, "SetSite", func(srv SiteControllerServiceServer, ctx context.Context, req *SetSiteRequest) (interface{}, error) {
			return srv.SetSite(ctx, req)
		})
	}},
	{MethodName: "SetService", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[SetServiceRequest](s, c, d, i, "SetService", func(srv SiteControllerServiceServer, ctx context.Context, req *SetServiceRequest) (interface{}, error) {
			return srv.SetService(ctx, req)
		})
	}},
	{MethodName: "SetRadio", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[SetRadioRequest](s, c, d, i, "SetRadio", func(srv SiteControllerServiceServer, ctx context.Context, req *SetRadioRequest) (interface{}, error) {
			return srv.SetRadio(ctx, req)
		})
	}},
	{MethodName: "GetSiteState", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[GetSiteStateRequest](s, c, d, i, "GetSiteState", func(srv SiteControllerServiceServer, ctx context.Context, req *GetSiteStateRequest) (interface{}, error) {
			return srv.GetSiteState(ctx, req)
		})
	}},
	{MethodName: "GetSwitchPolicy", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[GetSwitchPolicyRequest](s, c, d, i, "GetSwitchPolicy", func(srv SiteControllerServiceServer, ctx context.Context, req *GetSwitchPolicyRequest) (interface{}, error) {
			return srv.GetSwitchPolicy(ctx, req)
		})
	}},
	{MethodName: "RefreshSwitchPolicy", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[RefreshSwitchPolicyRequest](s, c, d, i, "RefreshSwitchPolicy", func(srv SiteControllerServiceServer, ctx context.Context, req *RefreshSwitchPolicyRequest) (interface{}, error) {
			return srv.RefreshSwitchPolicy(ctx, req)
		})
	}},
	{MethodName: "ReportSwitchPolicy", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[ReportSwitchPolicyRequest](s, c, d, i, "ReportSwitchPolicy", func(srv SiteControllerServiceServer, ctx context.Context, req *ReportSwitchPolicyRequest) (interface{}, error) {
			return srv.ReportSwitchPolicy(ctx, req)
		})
	}},
	{MethodName: "PowerCycleNode", Handler: func(s interface{}, c context.Context, d func(interface{}) error, i grpc.UnaryServerInterceptor) (interface{}, error) {
		return handler[PowerCycleNodeRequest](s, c, d, i, "PowerCycleNode", func(srv SiteControllerServiceServer, ctx context.Context, req *PowerCycleNodeRequest) (interface{}, error) {
			return srv.PowerCycleNode(ctx, req)
		})
	}},
}, Streams: []grpc.StreamDesc{}, Metadata: "site_controller.proto"}
