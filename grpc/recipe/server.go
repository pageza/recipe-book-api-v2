// File: grpc/recipe/server.go
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
		// For now, these additional fields are not used by the client.
		NutritionalInfo:   req.NutritionalInfo,
		AllergyDisclaimer: req.AllergyDisclaimer,
		Appliances:        req.Appliances,
	}

	if err := s.svc.CreateRecipe(recipe); err != nil {
		return nil, fmt.Errorf("failed to create recipe: %v", err)
	}

	return &pb.CreateRecipeResponse{
		RecipeId: recipe.ID,
		Message:  "Recipe created successfully",
	}, nil
}

// GetRecipe implements the GetRecipe RPC.
func (s *Server) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.GetRecipeResponse, error) {
	recipe, err := s.svc.GetRecipe(req.RecipeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %v", err)
	}

	resp := &pb.GetRecipeResponse{
		RecipeId:          recipe.ID,
		Title:             recipe.Title,
		Ingredients:       recipe.Ingredients,
		Steps:             recipe.Steps,
		NutritionalInfo:   recipe.NutritionalInfo,
		AllergyDisclaimer: recipe.AllergyDisclaimer,
		Appliances:        recipe.Appliances,
		CreatedAt:         recipe.CreatedAt.Unix(), // assuming Unix seconds; adjust if necessary
		UpdatedAt:         recipe.UpdatedAt.Unix(),
	}
	return resp, nil
}

// ListRecipes implements the ListRecipes RPC.
func (s *Server) ListRecipes(ctx context.Context, req *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	recipes, err := s.svc.GetAllRecipes()
	if err != nil {
		return nil, fmt.Errorf("failed to list recipes: %v", err)
	}

	var protoRecipes []*pb.GetRecipeResponse
	for _, r := range recipes {
		protoRecipes = append(protoRecipes, &pb.GetRecipeResponse{
			RecipeId:          r.ID,
			Title:             r.Title,
			Ingredients:       r.Ingredients,
			Steps:             r.Steps,
			NutritionalInfo:   r.NutritionalInfo,
			AllergyDisclaimer: r.AllergyDisclaimer,
			Appliances:        r.Appliances,
			CreatedAt:         r.CreatedAt.Unix(),
			UpdatedAt:         r.UpdatedAt.Unix(),
		})
	}

	return &pb.ListRecipesResponse{
		Recipes: protoRecipes,
	}, nil
}

// QueryRecipe implements the QueryRecipe RPC.
// This is a placeholder; you'll need to implement your RAG logic here later.
func (s *Server) QueryRecipe(ctx context.Context, req *pb.RecipeQueryRequest) (*pb.RecipeQueryResponse, error) {
	// For now, we'll simulate a query by retrieving the recipe by an ID that matches the query string.
	recipe, err := s.svc.GetRecipeByQuery(req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipe: %v", err)
	}

	// In a full implementation, you would also generate alternative suggestions.
	resp := &pb.RecipeQueryResponse{
		PrimaryRecipe: &pb.GetRecipeResponse{
			RecipeId:          recipe.ID,
			Title:             recipe.Title,
			Ingredients:       recipe.Ingredients,
			Steps:             recipe.Steps,
			NutritionalInfo:   recipe.NutritionalInfo,
			AllergyDisclaimer: recipe.AllergyDisclaimer,
			Appliances:        recipe.Appliances,
			CreatedAt:         recipe.CreatedAt.Unix(),
			UpdatedAt:         recipe.UpdatedAt.Unix(),
		},
		// Alternative recipes can be added here.
	}
	return resp, nil
}
