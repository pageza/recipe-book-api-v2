// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: recipe/recipe.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RecipeService_CreateRecipe_FullMethodName = "/recipe.RecipeService/CreateRecipe"
	RecipeService_GetRecipe_FullMethodName    = "/recipe.RecipeService/GetRecipe"
	RecipeService_ListRecipes_FullMethodName  = "/recipe.RecipeService/ListRecipes"
	RecipeService_QueryRecipe_FullMethodName  = "/recipe.RecipeService/QueryRecipe"
)

// RecipeServiceClient is the client API for RecipeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecipeServiceClient interface {
	CreateRecipe(ctx context.Context, in *CreateRecipeRequest, opts ...grpc.CallOption) (*CreateRecipeResponse, error)
	GetRecipe(ctx context.Context, in *GetRecipeRequest, opts ...grpc.CallOption) (*GetRecipeResponse, error)
	ListRecipes(ctx context.Context, in *ListRecipesRequest, opts ...grpc.CallOption) (*ListRecipesResponse, error)
	QueryRecipe(ctx context.Context, in *RecipeQueryRequest, opts ...grpc.CallOption) (*RecipeQueryResponse, error)
}

type recipeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecipeServiceClient(cc grpc.ClientConnInterface) RecipeServiceClient {
	return &recipeServiceClient{cc}
}

func (c *recipeServiceClient) CreateRecipe(ctx context.Context, in *CreateRecipeRequest, opts ...grpc.CallOption) (*CreateRecipeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRecipeResponse)
	err := c.cc.Invoke(ctx, RecipeService_CreateRecipe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipeServiceClient) GetRecipe(ctx context.Context, in *GetRecipeRequest, opts ...grpc.CallOption) (*GetRecipeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRecipeResponse)
	err := c.cc.Invoke(ctx, RecipeService_GetRecipe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipeServiceClient) ListRecipes(ctx context.Context, in *ListRecipesRequest, opts ...grpc.CallOption) (*ListRecipesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRecipesResponse)
	err := c.cc.Invoke(ctx, RecipeService_ListRecipes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipeServiceClient) QueryRecipe(ctx context.Context, in *RecipeQueryRequest, opts ...grpc.CallOption) (*RecipeQueryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecipeQueryResponse)
	err := c.cc.Invoke(ctx, RecipeService_QueryRecipe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecipeServiceServer is the server API for RecipeService service.
// All implementations must embed UnimplementedRecipeServiceServer
// for forward compatibility.
type RecipeServiceServer interface {
	CreateRecipe(context.Context, *CreateRecipeRequest) (*CreateRecipeResponse, error)
	GetRecipe(context.Context, *GetRecipeRequest) (*GetRecipeResponse, error)
	ListRecipes(context.Context, *ListRecipesRequest) (*ListRecipesResponse, error)
	QueryRecipe(context.Context, *RecipeQueryRequest) (*RecipeQueryResponse, error)
	mustEmbedUnimplementedRecipeServiceServer()
}

// UnimplementedRecipeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRecipeServiceServer struct{}

func (UnimplementedRecipeServiceServer) CreateRecipe(context.Context, *CreateRecipeRequest) (*CreateRecipeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRecipe not implemented")
}
func (UnimplementedRecipeServiceServer) GetRecipe(context.Context, *GetRecipeRequest) (*GetRecipeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecipe not implemented")
}
func (UnimplementedRecipeServiceServer) ListRecipes(context.Context, *ListRecipesRequest) (*ListRecipesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRecipes not implemented")
}
func (UnimplementedRecipeServiceServer) QueryRecipe(context.Context, *RecipeQueryRequest) (*RecipeQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRecipe not implemented")
}
func (UnimplementedRecipeServiceServer) mustEmbedUnimplementedRecipeServiceServer() {}
func (UnimplementedRecipeServiceServer) testEmbeddedByValue()                       {}

// UnsafeRecipeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecipeServiceServer will
// result in compilation errors.
type UnsafeRecipeServiceServer interface {
	mustEmbedUnimplementedRecipeServiceServer()
}

func RegisterRecipeServiceServer(s grpc.ServiceRegistrar, srv RecipeServiceServer) {
	// If the following call pancis, it indicates UnimplementedRecipeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RecipeService_ServiceDesc, srv)
}

func _RecipeService_CreateRecipe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRecipeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).CreateRecipe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_CreateRecipe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).CreateRecipe(ctx, req.(*CreateRecipeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipeService_GetRecipe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecipeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).GetRecipe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_GetRecipe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).GetRecipe(ctx, req.(*GetRecipeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipeService_ListRecipes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRecipesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).ListRecipes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_ListRecipes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).ListRecipes(ctx, req.(*ListRecipesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipeService_QueryRecipe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipeQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).QueryRecipe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_QueryRecipe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).QueryRecipe(ctx, req.(*RecipeQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecipeService_ServiceDesc is the grpc.ServiceDesc for RecipeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecipeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "recipe.RecipeService",
	HandlerType: (*RecipeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRecipe",
			Handler:    _RecipeService_CreateRecipe_Handler,
		},
		{
			MethodName: "GetRecipe",
			Handler:    _RecipeService_GetRecipe_Handler,
		},
		{
			MethodName: "ListRecipes",
			Handler:    _RecipeService_ListRecipes_Handler,
		},
		{
			MethodName: "QueryRecipe",
			Handler:    _RecipeService_QueryRecipe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recipe/recipe.proto",
}
