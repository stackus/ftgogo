// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package accountingpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccountingServiceClient is the client API for AccountingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountingServiceClient interface {
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	DisableAccount(ctx context.Context, in *DisableAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	EnableAccount(ctx context.Context, in *EnableAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type accountingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountingServiceClient(cc grpc.ClientConnInterface) AccountingServiceClient {
	return &accountingServiceClient{cc}
}

func (c *accountingServiceClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, "/accountingpb.AccountingService/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingServiceClient) DisableAccount(ctx context.Context, in *DisableAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/accountingpb.AccountingService/DisableAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountingServiceClient) EnableAccount(ctx context.Context, in *EnableAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/accountingpb.AccountingService/EnableAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountingServiceServer is the server API for AccountingService service.
// All implementations must embed UnimplementedAccountingServiceServer
// for forward compatibility
type AccountingServiceServer interface {
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	DisableAccount(context.Context, *DisableAccountRequest) (*emptypb.Empty, error)
	EnableAccount(context.Context, *EnableAccountRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAccountingServiceServer()
}

// UnimplementedAccountingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountingServiceServer struct {
}

func (UnimplementedAccountingServiceServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedAccountingServiceServer) DisableAccount(context.Context, *DisableAccountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableAccount not implemented")
}
func (UnimplementedAccountingServiceServer) EnableAccount(context.Context, *EnableAccountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableAccount not implemented")
}
func (UnimplementedAccountingServiceServer) mustEmbedUnimplementedAccountingServiceServer() {}

// UnsafeAccountingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountingServiceServer will
// result in compilation errors.
type UnsafeAccountingServiceServer interface {
	mustEmbedUnimplementedAccountingServiceServer()
}

func RegisterAccountingServiceServer(s grpc.ServiceRegistrar, srv AccountingServiceServer) {
	s.RegisterService(&AccountingService_ServiceDesc, srv)
}

func _AccountingService_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServiceServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accountingpb.AccountingService/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServiceServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingService_DisableAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServiceServer).DisableAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accountingpb.AccountingService/DisableAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServiceServer).DisableAccount(ctx, req.(*DisableAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountingService_EnableAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountingServiceServer).EnableAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accountingpb.AccountingService/EnableAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountingServiceServer).EnableAccount(ctx, req.(*EnableAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountingService_ServiceDesc is the grpc.ServiceDesc for AccountingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "accountingpb.AccountingService",
	HandlerType: (*AccountingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccount",
			Handler:    _AccountingService_GetAccount_Handler,
		},
		{
			MethodName: "DisableAccount",
			Handler:    _AccountingService_DisableAccount_Handler,
		},
		{
			MethodName: "EnableAccount",
			Handler:    _AccountingService_EnableAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "accountingapi/pb/service.proto",
}
