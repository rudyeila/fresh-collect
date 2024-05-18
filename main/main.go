package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/rudyeila/hello-fresh-go-client/config"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh/model"
)

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	log := slog.New(slog.NewJSONHandler(os.Stdout, options))

	cfg, err := config.New()
	if err != nil {
		log.Error("creating config", "error", err.Error())
	}

	hf := hellofresh.New(cfg, log)
	recipe, err := hf.GetRecipe("65d324a4833ed7a0e818b1b5")
	if err != nil {
		log.Error("getting recipe", "id", "64df2a4d614f75555c20edba", "error", err.Error())
	}

	shippedIngr := make([]model.Ingredient, 0)
	for _, ing := range recipe.Ingredients {
		if ing.Shipped {
			fmt.Println(ing.Name)
			shippedIngr = append(shippedIngr, ing)
		}
	}

	var yieldForTwo model.Yield
	for _, y := range recipe.Yields {
		if y.Yields == 2 {
			yieldForTwo = y
			break
		}
	}

	ingredientsWithAmount := make([]model.IngredientWithAmount, len(shippedIngr))
	for i, ingr := range shippedIngr {
		for _, yIngr := range yieldForTwo.Ingredients {
			if ingr.Id == yIngr.Id {
				ingredientsWithAmount[i] = model.IngredientWithAmount{
					Id:     ingr.Id,
					Uuid:   ingr.Uuid,
					Name:   ingr.Name,
					Amount: yIngr.Amount,
					Unit:   yIngr.Unit,
				}
			}
		}
	}

	err = WriteStructToJSONFile(ingredientsWithAmount, "ingredients.json")
	if err != nil {
		log.Error("Writing ingredients to file", "error", err)

	}
}

// WriteStructToJSONFile writes a given struct to a specified file as JSON
func WriteStructToJSONFile(v interface{}, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a JSON encoder and write the struct to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // for pretty printing
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("could not encode to JSON: %v", err)
	}

	return nil
}
