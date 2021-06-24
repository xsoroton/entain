// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sport

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SportsClient is the client API for Sports service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SportsClient interface {
	// GetSport by ID
	GetSport(ctx context.Context, in *GetSportRequest, opts ...grpc.CallOption) (*GetSportResponse, error)
	// ListSports returns a list of all sports.
	ListSports(ctx context.Context, in *ListSportsRequest, opts ...grpc.CallOption) (*ListSportsResponse, error)
}

type sportsClient struct {
	cc grpc.ClientConnInterface
}

func NewSportsClient(cc grpc.ClientConnInterface) SportsClient {
	return &sportsClient{cc}
}

func (c *sportsClient) GetSport(ctx context.Context, in *GetSportRequest, opts ...grpc.CallOption) (*GetSportResponse, error) {
	out := new(GetSportResponse)
	err := c.cc.Invoke(ctx, "/sport.Sports/GetSport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sportsClient) ListSports(ctx context.Context, in *ListSportsRequest, opts ...grpc.CallOption) (*ListSportsResponse, error) {
	out := new(ListSportsResponse)
	err := c.cc.Invoke(ctx, "/sport.Sports/ListSports", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SportsServer is the server API for Sports service.
// All implementations must embed UnimplementedSportsServer
// for forward compatibility
type SportsServer interface {
	// GetSport by ID
	GetSport(context.Context, *GetSportRequest) (*GetSportResponse, error)
	// ListSports returns a list of all sports.
	ListSports(context.Context, *ListSportsRequest) (*ListSportsResponse, error)
	mustEmbedUnimplementedSportsServer()
}

// UnimplementedSportsServer must be embedded to have forward compatible implementations.
type UnimplementedSportsServer struct {
}

func (UnimplementedSportsServer) GetSport(context.Context, *GetSportRequest) (*GetSportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSport not implemented")
}
func (UnimplementedSportsServer) ListSports(context.Context, *ListSportsRequest) (*ListSportsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSports not implemented")
}
func (UnimplementedSportsServer) mustEmbedUnimplementedSportsServer() {}

// UnsafeSportsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SportsServer will
// result in compilation errors.
type UnsafeSportsServer interface {
	mustEmbedUnimplementedSportsServer()
}

func RegisterSportsServer(s grpc.ServiceRegistrar, srv SportsServer) {
	s.RegisterService(&Sports_ServiceDesc, srv)
}

func _Sports_GetSport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportsServer).GetSport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sport.Sports/GetSport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportsServer).GetSport(ctx, req.(*GetSportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sports_ListSports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportsServer).ListSports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sport.Sports/ListSports",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportsServer).ListSports(ctx, req.(*ListSportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sports_ServiceDesc is the grpc.ServiceDesc for Sports service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sports_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sport.Sports",
	HandlerType: (*SportsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSport",
			Handler:    _Sports_GetSport_Handler,
		},
		{
			MethodName: "ListSports",
			Handler:    _Sports_ListSports_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sport/sport.proto",
}