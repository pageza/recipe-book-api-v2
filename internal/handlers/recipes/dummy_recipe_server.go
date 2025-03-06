// cursor--DummyRecipeServer: minimal gRPC server implementation for testing purposes.
package recipes

import (
	"context"

	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

// DummyRecipeServer implements pb.RecipeServiceServer with trivial methods.
type DummyRecipeServer struct {
	pb.UnimplementedRecipeServiceServer
}

// NewDummyRecipeServer creates a new instance of DummyRecipeServer.
func NewDummyRecipeServer() *DummyRecipeServer {
	return &DummyRecipeServer{}
}

// GetRecipe returns an empty GetRecipeResponse immediately.
func (s *DummyRecipeServer) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.GetRecipeResponse, error) {
	return &pb.GetRecipeResponse{}, nil
}

// QueryRecipe returns an empty RecipeQueryResponse immediately.
func (s *DummyRecipeServer) QueryRecipe(ctx context.Context, req *pb.RecipeQueryRequest) (*pb.RecipeQueryResponse, error) {
	return &pb.RecipeQueryResponse{}, nil
}
