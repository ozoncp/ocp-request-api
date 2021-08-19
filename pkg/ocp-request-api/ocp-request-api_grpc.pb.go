// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ocp_request_api

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

// OcpRequestApiClient is the client API for OcpRequestApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OcpRequestApiClient interface {
	// ListRequestV1 returns a list of user Requests.
	ListRequestV1(ctx context.Context, in *ListRequestsV1Request, opts ...grpc.CallOption) (*ListRequestsV1Response, error)
	// DescribeTaskV1 returns detailed information of a given Request.
	DescribeTaskV1(ctx context.Context, in *DescribeRequestV1Request, opts ...grpc.CallOption) (*DescribeTaskV1Response, error)
	// CreateRequestV1 creates new request. Returns id of created object.
	CreateRequestV1(ctx context.Context, in *CreateRequestV1Request, opts ...grpc.CallOption) (*CreateRequestV1Response, error)
	// RemoveRequestV1 removes user request by a its by.
	// Returns a bool flag indicating if object actually existed and hence removed.
	RemoveRequestV1(ctx context.Context, in *RemoveRequestV1Request, opts ...grpc.CallOption) (*RemoveRequestV1Response, error)
}

type ocpRequestApiClient struct {
	cc grpc.ClientConnInterface
}

func NewOcpRequestApiClient(cc grpc.ClientConnInterface) OcpRequestApiClient {
	return &ocpRequestApiClient{cc}
}

func (c *ocpRequestApiClient) ListRequestV1(ctx context.Context, in *ListRequestsV1Request, opts ...grpc.CallOption) (*ListRequestsV1Response, error) {
	out := new(ListRequestsV1Response)
	err := c.cc.Invoke(ctx, "/ocp.task.api.OcpRequestApi/ListRequestV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ocpRequestApiClient) DescribeTaskV1(ctx context.Context, in *DescribeRequestV1Request, opts ...grpc.CallOption) (*DescribeTaskV1Response, error) {
	out := new(DescribeTaskV1Response)
	err := c.cc.Invoke(ctx, "/ocp.task.api.OcpRequestApi/DescribeTaskV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ocpRequestApiClient) CreateRequestV1(ctx context.Context, in *CreateRequestV1Request, opts ...grpc.CallOption) (*CreateRequestV1Response, error) {
	out := new(CreateRequestV1Response)
	err := c.cc.Invoke(ctx, "/ocp.task.api.OcpRequestApi/CreateRequestV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ocpRequestApiClient) RemoveRequestV1(ctx context.Context, in *RemoveRequestV1Request, opts ...grpc.CallOption) (*RemoveRequestV1Response, error) {
	out := new(RemoveRequestV1Response)
	err := c.cc.Invoke(ctx, "/ocp.task.api.OcpRequestApi/RemoveRequestV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OcpRequestApiServer is the server API for OcpRequestApi service.
// All implementations must embed UnimplementedOcpRequestApiServer
// for forward compatibility
type OcpRequestApiServer interface {
	// ListRequestV1 returns a list of user Requests.
	ListRequestV1(context.Context, *ListRequestsV1Request) (*ListRequestsV1Response, error)
	// DescribeTaskV1 returns detailed information of a given Request.
	DescribeTaskV1(context.Context, *DescribeRequestV1Request) (*DescribeTaskV1Response, error)
	// CreateRequestV1 creates new request. Returns id of created object.
	CreateRequestV1(context.Context, *CreateRequestV1Request) (*CreateRequestV1Response, error)
	// RemoveRequestV1 removes user request by a its by.
	// Returns a bool flag indicating if object actually existed and hence removed.
	RemoveRequestV1(context.Context, *RemoveRequestV1Request) (*RemoveRequestV1Response, error)
	mustEmbedUnimplementedOcpRequestApiServer()
}

// UnimplementedOcpRequestApiServer must be embedded to have forward compatible implementations.
type UnimplementedOcpRequestApiServer struct {
}

func (UnimplementedOcpRequestApiServer) ListRequestV1(context.Context, *ListRequestsV1Request) (*ListRequestsV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRequestV1 not implemented")
}
func (UnimplementedOcpRequestApiServer) DescribeTaskV1(context.Context, *DescribeRequestV1Request) (*DescribeTaskV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeTaskV1 not implemented")
}
func (UnimplementedOcpRequestApiServer) CreateRequestV1(context.Context, *CreateRequestV1Request) (*CreateRequestV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRequestV1 not implemented")
}
func (UnimplementedOcpRequestApiServer) RemoveRequestV1(context.Context, *RemoveRequestV1Request) (*RemoveRequestV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveRequestV1 not implemented")
}
func (UnimplementedOcpRequestApiServer) mustEmbedUnimplementedOcpRequestApiServer() {}

// UnsafeOcpRequestApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OcpRequestApiServer will
// result in compilation errors.
type UnsafeOcpRequestApiServer interface {
	mustEmbedUnimplementedOcpRequestApiServer()
}

func RegisterOcpRequestApiServer(s grpc.ServiceRegistrar, srv OcpRequestApiServer) {
	s.RegisterService(&OcpRequestApi_ServiceDesc, srv)
}

func _OcpRequestApi_ListRequestV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequestsV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OcpRequestApiServer).ListRequestV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocp.task.api.OcpRequestApi/ListRequestV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OcpRequestApiServer).ListRequestV1(ctx, req.(*ListRequestsV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _OcpRequestApi_DescribeTaskV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeRequestV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OcpRequestApiServer).DescribeTaskV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocp.task.api.OcpRequestApi/DescribeTaskV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OcpRequestApiServer).DescribeTaskV1(ctx, req.(*DescribeRequestV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _OcpRequestApi_CreateRequestV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequestV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OcpRequestApiServer).CreateRequestV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocp.task.api.OcpRequestApi/CreateRequestV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OcpRequestApiServer).CreateRequestV1(ctx, req.(*CreateRequestV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _OcpRequestApi_RemoveRequestV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequestV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OcpRequestApiServer).RemoveRequestV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocp.task.api.OcpRequestApi/RemoveRequestV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OcpRequestApiServer).RemoveRequestV1(ctx, req.(*RemoveRequestV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// OcpRequestApi_ServiceDesc is the grpc.ServiceDesc for OcpRequestApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OcpRequestApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ocp.task.api.OcpRequestApi",
	HandlerType: (*OcpRequestApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRequestV1",
			Handler:    _OcpRequestApi_ListRequestV1_Handler,
		},
		{
			MethodName: "DescribeTaskV1",
			Handler:    _OcpRequestApi_DescribeTaskV1_Handler,
		},
		{
			MethodName: "CreateRequestV1",
			Handler:    _OcpRequestApi_CreateRequestV1_Handler,
		},
		{
			MethodName: "RemoveRequestV1",
			Handler:    _OcpRequestApi_RemoveRequestV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ocp-request-api/ocp-request-api.proto",
}
