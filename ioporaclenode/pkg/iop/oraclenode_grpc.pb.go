// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package iop

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// OracleNodeClient is the client API for OracleNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OracleNodeClient interface {
	ValidateTransaction(ctx context.Context, in *ValidateTransactionRequest, opts ...grpc.CallOption) (*ValidateTransactionResponse, error)
	ProcessDeal(ctx context.Context, in *ProcessDealRequest, opts ...grpc.CallOption) (*ProcessDealResponse, error)
	ProcessResponse(ctx context.Context, in *ProcessResponseRequest, opts ...grpc.CallOption) (*ProcessResponseResponse, error)
	ProcessJustification(ctx context.Context, in *ProcessJustificationRequest, opts ...grpc.CallOption) (*ProcessJustificationResponse, error)
}

type oracleNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewOracleNodeClient(cc grpc.ClientConnInterface) OracleNodeClient {
	return &oracleNodeClient{cc}
}

func (c *oracleNodeClient) ValidateTransaction(ctx context.Context, in *ValidateTransactionRequest, opts ...grpc.CallOption) (*ValidateTransactionResponse, error) {
	out := new(ValidateTransactionResponse)
	err := c.cc.Invoke(ctx, "/iop.OracleNode/ValidateTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oracleNodeClient) ProcessDeal(ctx context.Context, in *ProcessDealRequest, opts ...grpc.CallOption) (*ProcessDealResponse, error) {
	out := new(ProcessDealResponse)
	err := c.cc.Invoke(ctx, "/iop.OracleNode/ProcessDeal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oracleNodeClient) ProcessResponse(ctx context.Context, in *ProcessResponseRequest, opts ...grpc.CallOption) (*ProcessResponseResponse, error) {
	out := new(ProcessResponseResponse)
	err := c.cc.Invoke(ctx, "/iop.OracleNode/ProcessResponse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oracleNodeClient) ProcessJustification(ctx context.Context, in *ProcessJustificationRequest, opts ...grpc.CallOption) (*ProcessJustificationResponse, error) {
	out := new(ProcessJustificationResponse)
	err := c.cc.Invoke(ctx, "/iop.OracleNode/ProcessJustification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OracleNodeServer is the server API for OracleNode service.
// All implementations must embed UnimplementedOracleNodeServer
// for forward compatibility
type OracleNodeServer interface {
	ValidateTransaction(context.Context, *ValidateTransactionRequest) (*ValidateTransactionResponse, error)
	ProcessDeal(context.Context, *ProcessDealRequest) (*ProcessDealResponse, error)
	ProcessResponse(context.Context, *ProcessResponseRequest) (*ProcessResponseResponse, error)
	ProcessJustification(context.Context, *ProcessJustificationRequest) (*ProcessJustificationResponse, error)
	mustEmbedUnimplementedOracleNodeServer()
}

// UnimplementedOracleNodeServer must be embedded to have forward compatible implementations.
type UnimplementedOracleNodeServer struct {
}

func (UnimplementedOracleNodeServer) ValidateTransaction(context.Context, *ValidateTransactionRequest) (*ValidateTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateTransaction not implemented")
}
func (UnimplementedOracleNodeServer) ProcessDeal(context.Context, *ProcessDealRequest) (*ProcessDealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessDeal not implemented")
}
func (UnimplementedOracleNodeServer) ProcessResponse(context.Context, *ProcessResponseRequest) (*ProcessResponseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessResponse not implemented")
}
func (UnimplementedOracleNodeServer) ProcessJustification(context.Context, *ProcessJustificationRequest) (*ProcessJustificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessJustification not implemented")
}
func (UnimplementedOracleNodeServer) mustEmbedUnimplementedOracleNodeServer() {}

// UnsafeOracleNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OracleNodeServer will
// result in compilation errors.
type UnsafeOracleNodeServer interface {
	mustEmbedUnimplementedOracleNodeServer()
}

func RegisterOracleNodeServer(s grpc.ServiceRegistrar, srv OracleNodeServer) {
	s.RegisterService(&OracleNode_ServiceDesc, srv)
}

func _OracleNode_ValidateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OracleNodeServer).ValidateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iop.OracleNode/ValidateTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OracleNodeServer).ValidateTransaction(ctx, req.(*ValidateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OracleNode_ProcessDeal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessDealRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OracleNodeServer).ProcessDeal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iop.OracleNode/ProcessDeal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OracleNodeServer).ProcessDeal(ctx, req.(*ProcessDealRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OracleNode_ProcessResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessResponseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OracleNodeServer).ProcessResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iop.OracleNode/ProcessResponse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OracleNodeServer).ProcessResponse(ctx, req.(*ProcessResponseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OracleNode_ProcessJustification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessJustificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OracleNodeServer).ProcessJustification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iop.OracleNode/ProcessJustification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OracleNodeServer).ProcessJustification(ctx, req.(*ProcessJustificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OracleNode_ServiceDesc is the grpc.ServiceDesc for OracleNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OracleNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iop.OracleNode",
	HandlerType: (*OracleNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateTransaction",
			Handler:    _OracleNode_ValidateTransaction_Handler,
		},
		{
			MethodName: "ProcessDeal",
			Handler:    _OracleNode_ProcessDeal_Handler,
		},
		{
			MethodName: "ProcessResponse",
			Handler:    _OracleNode_ProcessResponse_Handler,
		},
		{
			MethodName: "ProcessJustification",
			Handler:    _OracleNode_ProcessJustification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oraclenode.proto",
}
