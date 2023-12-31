// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: proto_ventas/ventas.proto

package ventas

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

// OrderProcessingClient is the client API for OrderProcessing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderProcessingClient interface {
	ProcessOrder(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
}

type orderProcessingClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderProcessingClient(cc grpc.ClientConnInterface) OrderProcessingClient {
	return &orderProcessingClient{cc}
}

func (c *orderProcessingClient) ProcessOrder(ctx context.Context, in *OrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, "/ventas.OrderProcessing/ProcessOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderProcessingServer is the server API for OrderProcessing service.
// All implementations must embed UnimplementedOrderProcessingServer
// for forward compatibility
type OrderProcessingServer interface {
	ProcessOrder(context.Context, *OrderRequest) (*OrderResponse, error)
	mustEmbedUnimplementedOrderProcessingServer()
}

// UnimplementedOrderProcessingServer must be embedded to have forward compatible implementations.
type UnimplementedOrderProcessingServer struct {
}

func (UnimplementedOrderProcessingServer) ProcessOrder(context.Context, *OrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessOrder not implemented")
}
func (UnimplementedOrderProcessingServer) mustEmbedUnimplementedOrderProcessingServer() {}

// UnsafeOrderProcessingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderProcessingServer will
// result in compilation errors.
type UnsafeOrderProcessingServer interface {
	mustEmbedUnimplementedOrderProcessingServer()
}

func RegisterOrderProcessingServer(s grpc.ServiceRegistrar, srv OrderProcessingServer) {
	s.RegisterService(&OrderProcessing_ServiceDesc, srv)
}

func _OrderProcessing_ProcessOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderProcessingServer).ProcessOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ventas.OrderProcessing/ProcessOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderProcessingServer).ProcessOrder(ctx, req.(*OrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderProcessing_ServiceDesc is the grpc.ServiceDesc for OrderProcessing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderProcessing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ventas.OrderProcessing",
	HandlerType: (*OrderProcessingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessOrder",
			Handler:    _OrderProcessing_ProcessOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto_ventas/ventas.proto",
}
