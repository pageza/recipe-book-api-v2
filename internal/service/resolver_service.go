package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pageza/recipe-book-api-v2/internal/models"
)

// ResolverServiceURL is the endpoint for the resolver microservice.
// In production, this is typically loaded from configuration.
var ResolverServiceURL = "http://localhost:8081/resolve"

// ResolutionRequest defines the request payload for recipe resolution.
type ResolutionRequest struct {
	Query string `json:"query"`
}

// ResolutionResponse defines the expected response payload.
type ResolutionResponse struct {
	PrimaryRecipe      *models.Recipe   `json:"primary_recipe"`
	AlternativeRecipes []*models.Recipe `json:"alternative_recipes"`
}

// CallResolver sends a request to the external resolver microservice and returns its response.
func CallResolver(query string) (*ResolutionResponse, error) {
	reqPayload := ResolutionRequest{Query: query}
	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ResolverServiceURL, "application/json", bytes.NewBuffer(reqBytes))
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
