package recipe

import (
	"context"
	"fmt"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
)

// Server implements the gRPC RecipeService.
type Server struct {
	pb.UnimplementedRecipeServiceServer
	svc service.RecipeService
}

// NewServer creates a new Recipe gRPC server.
func NewServer(svc service.RecipeService) *Server {
	return &Server{svc: svc}
}

// CreateRecipe implements the CreateRecipe RPC.
func (s *Server) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.CreateRecipeResponse, error) {
	recipe := &models.Recipe{
		Title:       req.Title,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
	}

	if err := s.svc.CreateRecipe(recipe); err != nil {
		return nil, fmt.Errorf("failed to create recipe: %v", err)
	}

	return &pb.CreateRecipeResponse{
		RecipeId: recipe.ID,
		Message:  "Recipe created successfully",
	}, nil
}
