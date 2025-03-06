// File: grpc/recipe/server.go
package recipe

//nolint:unusedwrite // false positive: field assignments are used in the gRPC response
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

// GetRecipe implements the GetRecipe RPC.
// It retrieves a recipe by its ID from the service layer.
func (s *Server) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.GetRecipeResponse, error) {
	recipe, err := s.svc.GetRecipe(req.RecipeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %v", err)
	}

	// Convert the internal Recipe model into a gRPC response.
	resp := &pb.GetRecipeResponse{
		RecipeId:          recipe.ID,
		Title:             recipe.Title,
		Ingredients:       recipe.Ingredients,
		Steps:             recipe.Steps,
		NutritionalInfo:   recipe.NutritionalInfo,
		AllergyDisclaimer: recipe.AllergyDisclaimer,
		Appliances:        recipe.Appliances,
		CreatedAt:         recipe.CreatedAt.Unix(),
		UpdatedAt:         recipe.UpdatedAt.Unix(),
	}
	return resp, nil
}

// QueryRecipe implements the QueryRecipe RPC.
// This endpoint handles both list (when query is empty) and advanced search queries.
func (s *Server) QueryRecipe(ctx context.Context, req *pb.RecipeQueryRequest) (*pb.RecipeQueryResponse, error) {
	// Use getters provided by the generated proto code.
	queryReq := &models.RecipeQueryRequest{
		Query:  req.Query,
		UserID: req.UserId,
		Filter: req.Filter,
		Page:   int(req.Page),
		Limit:  int(req.Limit),
	}

	// Delegate query processing to the service layer.
	queryResp, err := s.svc.QueryRecipes(queryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %v", err)
	}

	// Convert each internal Recipe into a proto GetRecipeResponse.
	var protoRecipes []*pb.GetRecipeResponse
	for _, r := range queryResp.Recipes {
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

	// Construct the final gRPC response including pagination details.
	resp := new(pb.RecipeQueryResponse)
	resp.Recipes = protoRecipes
	resp.Page = int32(queryResp.Page)
	resp.Limit = int32(queryResp.Limit)
	resp.Total = int32(queryResp.Total)
	return resp, nil
}
