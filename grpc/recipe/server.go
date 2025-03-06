// File: grpc/recipe/server.go
package recipe

//nolint:unusedwrite // false positive: field assignments are used in the gRPC response
import (
	"context"
	"fmt"
	"strings"

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
// It retrieves a recipe by its ID and converts the internal model into a gRPC response.
func (s *Server) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.GetRecipeResponse, error) {
	recipe, err := s.svc.GetRecipe(req.RecipeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %v", err)
	}

	ingredientsStr := strings.Join(recipe.Ingredients, ", ")
	stepsStr := strings.Join(recipe.Steps, ", ")
	appliancesStr := strings.Join(recipe.Appliances, ", ")
	nutritionalInfoStr := fmt.Sprintf("Calories: %.0f, Protein: %.0f, Carbs: %.0f, Fat: %.0f, Fiber: %.0f",
		recipe.NutritionalInfo.Calories,
		recipe.NutritionalInfo.Protein,
		recipe.NutritionalInfo.Carbohydrates,
		recipe.NutritionalInfo.Fat,
		recipe.NutritionalInfo.Fiber)

	resp := &pb.GetRecipeResponse{
		RecipeId:          recipe.ID,
		Title:             recipe.Title,
		Ingredients:       ingredientsStr,
		Steps:             stepsStr,
		NutritionalInfo:   nutritionalInfoStr,
		AllergyDisclaimer: recipe.AllergyDisclaimer,
		Appliances:        appliancesStr,
		CreatedAt:         recipe.CreatedAt.Unix(),
		UpdatedAt:         recipe.UpdatedAt.Unix(),
	}
	return resp, nil
}

// QueryRecipe implements the QueryRecipe RPC.
// This endpoint handles both list (when query is empty) and advanced search queries.
func (s *Server) QueryRecipe(ctx context.Context, req *pb.RecipeQueryRequest) (*pb.RecipeQueryResponse, error) {
	// Convert incoming proto request into an internal RecipeQueryRequest.
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

	var protoRecipes []*pb.GetRecipeResponse
	for _, rec := range queryResp.Recipes {
		// Convert slice fields into a comma-separated string
		ingredientsStr := strings.Join(rec.Ingredients, ", ")
		stepsStr := strings.Join(rec.Steps, ", ")
		appliancesStr := strings.Join(rec.Appliances, ", ")
		// Format the NutritionalInfo struct into a single string.
		nutritionalInfoStr := fmt.Sprintf("Calories: %.0f, Protein: %.0f, Carbs: %.0f, Fat: %.0f, Fiber: %.0f",
			rec.NutritionalInfo.Calories,
			rec.NutritionalInfo.Protein,
			rec.NutritionalInfo.Carbohydrates,
			rec.NutritionalInfo.Fat,
			rec.NutritionalInfo.Fiber)

		protoRecipes = append(protoRecipes, &pb.GetRecipeResponse{
			RecipeId:          rec.ID,
			Title:             rec.Title,
			Ingredients:       ingredientsStr,
			Steps:             stepsStr,
			NutritionalInfo:   nutritionalInfoStr,
			AllergyDisclaimer: rec.AllergyDisclaimer,
			Appliances:        appliancesStr,
			CreatedAt:         rec.CreatedAt.Unix(),
			UpdatedAt:         rec.UpdatedAt.Unix(),
		})
	}

	return &pb.RecipeQueryResponse{
		Recipes: protoRecipes,
		Page:    int32(queryResp.Page),
		Limit:   int32(queryResp.Limit),
		Total:   int32(queryResp.Total),
	}, nil
}
