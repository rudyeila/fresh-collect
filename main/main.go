package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/rudyeila/go-bring-api/bring"
	"github.com/rudyeila/hello-fresh-go-client/config"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh/client"
	"github.com/rudyeila/hello-fresh-go-client/hellofresh/service"
)

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	log := slog.New(slog.NewJSONHandler(os.Stdout, options))

	cfg, err := config.New()
	if err != nil {
		log.Error("creating config", "error", err.Error())
	}

	hf := client.NewClient(cfg.HelloFresh, log)
	b := bring.New(cfg.Bring, log)
	err = b.Login()
	if err != nil {
		log.Error("authentication user with Bring", "error", err)
	}

	svc := service.Service{Bring: b, HF: hf, Log: log}

	recIds := []string{"64df2a4d614f75555c20edba", "586250316121bb04b97342c2", "58343e5dd4d92c5781367e02"}
	ingredients, err := svc.GetMergedIngredients(true, recIds...)
	if err != nil {
		log.Error("Getting ingredients", "error", err)
	}

	err = WriteStructToJSONFile(ingredients, "ingredients.json")
	if err != nil {
		log.Error("Writing ingredients to file", "error", err)
	}

	err = svc.AddToShoppingList("HelloFresh", ingredients)
	if err != nil {
		log.Error("Adding ingredients to Bring", "error", err)
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

func ReadJSONFromFile(filename string, v interface{}) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the file's content
	err = json.NewDecoder(file).Decode(&v)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return nil
}
