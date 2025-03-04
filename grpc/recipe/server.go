// File: grpc/recipe/server.go
package recipe

import (
	"context"
	"fmt"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"github.com/pageza/recipe-book-api-v2/internal/service"
	pb "github.com/pageza/recipe-book-api-v2/proto/proto"
	"go.uber.org/zap"
)

// Server implements the RecipeService defined in the proto file.
type Server struct {
	svc service.RecipeService
	pb.UnimplementedRecipeServiceServer
}

// NewServer creates a new gRPC recipe server.
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
	recipes, err := s.svc.ListRecipes()
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

// QueryRecipe processes a recipe query using the full resolution flow.
func (s *Server) QueryRecipe(ctx context.Context, req *pb.RecipeQueryRequest) (*pb.RecipeQueryResponse, error) {
	zap.L().Info("Received QueryRecipe request", zap.String("query", req.Query))

	// Use the full resolution logic from the service layer.
	queryResp, err := s.svc.ResolveRecipeQuery(req.Query)
	if err != nil {
		zap.L().Warn("Failed to resolve recipe query", zap.Error(err))
		return nil, fmt.Errorf("failed to resolve recipe query: %v", err)
	}

	if len(queryResp.Recipes) == 0 {
		return nil, fmt.Errorf("no recipe found")
	}

	// Map the first returned recipe as the primary recipe.
	primary := mapRecipeToProto(queryResp.Recipes[0])

	resp := &pb.RecipeQueryResponse{
		PrimaryRecipe:      primary,
		AlternativeRecipes: []*pb.GetRecipeResponse{}, // Placeholder for alternatives.
	}

	return resp, nil
}

// mapRecipeToProto converts a Recipe model into its corresponding gRPC message.
func mapRecipeToProto(r *models.Recipe) *pb.GetRecipeResponse {
	return &pb.GetRecipeResponse{
		RecipeId:          r.ID,
		Title:             r.Title,
		Ingredients:       r.Ingredients,
		Steps:             r.Steps,
		NutritionalInfo:   r.NutritionalInfo,
		AllergyDisclaimer: r.AllergyDisclaimer,
		Appliances:        r.Appliances,
		CreatedAt:         r.CreatedAt.Unix(),
		UpdatedAt:         r.UpdatedAt.Unix(),
	}
}
