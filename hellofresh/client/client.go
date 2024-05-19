package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rudyeila/hello-fresh-go-client/hellofresh"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh/client/model"
)

type HelloFresh struct {
	client  *http.Client
	baseURL string
	log     *slog.Logger
	cfg     hellofresh.Config
}

func NewClient(cfg hellofresh.Config, logger *slog.Logger) *HelloFresh {
	client := &http.Client{
		Timeout: cfg.DefaultTimeout,
	}

	return &HelloFresh{
		cfg:     cfg,
		baseURL: cfg.BaseURL,
		client:  client,
		log:     logger,
	}
}

func (c *HelloFresh) GetRecipe(recipeID string) (*model.GetRecipeResponse, error) {
	url := fmt.Sprintf("%s/recipes/%s", c.baseURL, recipeID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.AccessToken))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	recipe := &model.GetRecipeResponse{}
	err = json.NewDecoder(res.Body).Decode(recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}
