package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pageza/recipe-book-api-v2/internal/models"
)

// These request and response types define the API contract with the resolver service.
type ResolutionRequest struct {
	Query string `json:"query"`
}

type ResolutionResponse struct {
	PrimaryRecipe      *models.Recipe   `json:"primary_recipe"`
	AlternativeRecipes []*models.Recipe `json:"alternative_recipes"`
}

// In production this URL should be read from the configuration.
var resolverServiceURL = "http://localhost:8081/resolve"

// ResolveRecipeQuery sends the recipe query to the resolver microservice.
func ResolveRecipeQuery(query string) (*ResolutionResponse, error) {
	reqPayload := ResolutionRequest{Query: query}
	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(resolverServiceURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resolver service returned status code: %d", resp.StatusCode)
	}

	var resolutionResp ResolutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&resolutionResp); err != nil {
		return nil, err
	}
	return &resolutionResp, nil
}
