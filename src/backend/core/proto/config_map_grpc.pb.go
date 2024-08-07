// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: config_map.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConfigMapClient is the client API for ConfigMap service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigMapClient interface {
	GetList(ctx context.Context, in *ConfigMapList, opts ...grpc.CallOption) (*ConfigMapListResponse, error)
	Save(ctx context.Context, in *ConfigMapSave, opts ...grpc.CallOption) (*empty.Empty, error)
	Delete(ctx context.Context, in *ConfigMapDelete, opts ...grpc.CallOption) (*empty.Empty, error)
}

type configMapClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigMapClient(cc grpc.ClientConnInterface) ConfigMapClient {
	return &configMapClient{cc}
}

func (c *configMapClient) GetList(ctx context.Context, in *ConfigMapList, opts ...grpc.CallOption) (*ConfigMapListResponse, error) {
	out := new(ConfigMapListResponse)
	err := c.cc.Invoke(ctx, "/com.cognix.ConfigMap/GetList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configMapClient) Save(ctx context.Context, in *ConfigMapSave, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.cognix.ConfigMap/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configMapClient) Delete(ctx context.Context, in *ConfigMapDelete, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/com.cognix.ConfigMap/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigMapServer is the server API for ConfigMap service.
// All implementations must embed UnimplementedConfigMapServer
// for forward compatibility
type ConfigMapServer interface {
	GetList(context.Context, *ConfigMapList) (*ConfigMapListResponse, error)
	Save(context.Context, *ConfigMapSave) (*empty.Empty, error)
	Delete(context.Context, *ConfigMapDelete) (*empty.Empty, error)
	mustEmbedUnimplementedConfigMapServer()
}

// UnimplementedConfigMapServer must be embedded to have forward compatible implementations.
type UnimplementedConfigMapServer struct {
}

func (UnimplementedConfigMapServer) GetList(context.Context, *ConfigMapList) (*ConfigMapListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedConfigMapServer) Save(context.Context, *ConfigMapSave) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedConfigMapServer) Delete(context.Context, *ConfigMapDelete) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedConfigMapServer) mustEmbedUnimplementedConfigMapServer() {}

// UnsafeConfigMapServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigMapServer will
// result in compilation errors.
type UnsafeConfigMapServer interface {
	mustEmbedUnimplementedConfigMapServer()
}

func RegisterConfigMapServer(s grpc.ServiceRegistrar, srv ConfigMapServer) {
	s.RegisterService(&ConfigMap_ServiceDesc, srv)
}

func _ConfigMap_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigMapList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigMapServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.cognix.ConfigMap/GetList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigMapServer).GetList(ctx, req.(*ConfigMapList))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigMap_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigMapSave)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigMapServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.cognix.ConfigMap/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigMapServer).Save(ctx, req.(*ConfigMapSave))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigMap_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigMapDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigMapServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.cognix.ConfigMap/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigMapServer).Delete(ctx, req.(*ConfigMapDelete))
	}
	return interceptor(ctx, in, info, handler)
}

// ConfigMap_ServiceDesc is the grpc.ServiceDesc for ConfigMap service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigMap_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.cognix.ConfigMap",
	HandlerType: (*ConfigMapServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetList",
			Handler:    _ConfigMap_GetList_Handler,
		},
		{
			MethodName: "Save",
			Handler:    _ConfigMap_Save_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ConfigMap_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config_map.proto",
}
