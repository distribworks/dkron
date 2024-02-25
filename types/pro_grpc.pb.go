// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: pro.proto

package types

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

// DkronProClient is the client API for DkronPro service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DkronProClient interface {
	ACLPolicyUpsert(ctx context.Context, in *ACLPolicyUpsertRequest, opts ...grpc.CallOption) (*ACLPolicyUpsertResponse, error)
	ACLPolicyDelete(ctx context.Context, in *ACLPolicyDeleteRequest, opts ...grpc.CallOption) (*ACLPolicyDeleteResponse, error)
}

type dkronProClient struct {
	cc grpc.ClientConnInterface
}

func NewDkronProClient(cc grpc.ClientConnInterface) DkronProClient {
	return &dkronProClient{cc}
}

func (c *dkronProClient) ACLPolicyUpsert(ctx context.Context, in *ACLPolicyUpsertRequest, opts ...grpc.CallOption) (*ACLPolicyUpsertResponse, error) {
	out := new(ACLPolicyUpsertResponse)
	err := c.cc.Invoke(ctx, "/types.DkronPro/ACLPolicyUpsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dkronProClient) ACLPolicyDelete(ctx context.Context, in *ACLPolicyDeleteRequest, opts ...grpc.CallOption) (*ACLPolicyDeleteResponse, error) {
	out := new(ACLPolicyDeleteResponse)
	err := c.cc.Invoke(ctx, "/types.DkronPro/ACLPolicyDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DkronProServer is the server API for DkronPro service.
// All implementations must embed UnimplementedDkronProServer
// for forward compatibility
type DkronProServer interface {
	ACLPolicyUpsert(context.Context, *ACLPolicyUpsertRequest) (*ACLPolicyUpsertResponse, error)
	ACLPolicyDelete(context.Context, *ACLPolicyDeleteRequest) (*ACLPolicyDeleteResponse, error)
	mustEmbedUnimplementedDkronProServer()
}

// UnimplementedDkronProServer must be embedded to have forward compatible implementations.
type UnimplementedDkronProServer struct {
}

func (UnimplementedDkronProServer) ACLPolicyUpsert(context.Context, *ACLPolicyUpsertRequest) (*ACLPolicyUpsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ACLPolicyUpsert not implemented")
}
func (UnimplementedDkronProServer) ACLPolicyDelete(context.Context, *ACLPolicyDeleteRequest) (*ACLPolicyDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ACLPolicyDelete not implemented")
}
func (UnimplementedDkronProServer) mustEmbedUnimplementedDkronProServer() {}

// UnsafeDkronProServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DkronProServer will
// result in compilation errors.
type UnsafeDkronProServer interface {
	mustEmbedUnimplementedDkronProServer()
}

func RegisterDkronProServer(s grpc.ServiceRegistrar, srv DkronProServer) {
	s.RegisterService(&DkronPro_ServiceDesc, srv)
}

func _DkronPro_ACLPolicyUpsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ACLPolicyUpsertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DkronProServer).ACLPolicyUpsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.DkronPro/ACLPolicyUpsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DkronProServer).ACLPolicyUpsert(ctx, req.(*ACLPolicyUpsertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DkronPro_ACLPolicyDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ACLPolicyDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DkronProServer).ACLPolicyDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.DkronPro/ACLPolicyDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DkronProServer).ACLPolicyDelete(ctx, req.(*ACLPolicyDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DkronPro_ServiceDesc is the grpc.ServiceDesc for DkronPro service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DkronPro_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "types.DkronPro",
	HandlerType: (*DkronProServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ACLPolicyUpsert",
			Handler:    _DkronPro_ACLPolicyUpsert_Handler,
		},
		{
			MethodName: "ACLPolicyDelete",
			Handler:    _DkronPro_ACLPolicyDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pro.proto",
}