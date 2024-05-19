package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/rudyeila/fresh-collect/config"
	"github.com/rudyeila/fresh-collect/hellofresh/client"
	"github.com/rudyeila/fresh-collect/hellofresh/service"
	"github.com/rudyeila/fresh-collect/hellofresh/service/model"
	"github.com/rudyeila/go-bring-api/bring"
	"github.com/urfave/cli/v2"
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

	app := &cli.App{
		Name:        "fresh-collect",
		Description: "a tool that helps plan your shopping list based on recipes",
		Authors: []*cli.Author{
			{
				Name:  "Rudy Ailabouni",
				Email: "eilabouni.rudy@gmail.com",
			},
		},
		Commands: []*cli.Command{
			GetParseRecipesCommand(svc),
			GetAddToListCommand(svc),
		},
	}

	if err = app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

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

func GetParseRecipesCommand(svc service.Service) *cli.Command {
	return &cli.Command{
		Name:    "parse",
		Aliases: []string{"p"},
		Usage:   "Parse recipes by IDs and get grouped ingredients",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "recipes",
				Aliases: []string{"r"},
				Usage:   "HelloFresh Recipe IDs to parse",
			},
			&cli.BoolFlag{
				Name:    "shipped",
				Aliases: []string{"s"},
				Value:   true,
				Usage:   "If true, only shipped ingredients are considered",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "File to store parsed and merged json in",
			},
		},
		Action: func(c *cli.Context) error {
			shipped := c.Bool("shipped")
			output := c.String("output")
			recIds := c.StringSlice("recipes")
			fmt.Println(shipped)
			fmt.Println(output)
			fmt.Println(recIds)

			ingredients, err := svc.GetMergedIngredients(shipped, recIds...)
			if err != nil {
				return err
			}

			if output != "" {
				err = WriteStructToJSONFile(ingredients, output)
				if err != nil {
					return err
				}
			} else {
				fmt.Println(ingredients)
			}

			return nil
		},
	}
}

func GetAddToListCommand(svc service.Service) *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add ingredients to shopping list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Required: true,
				Usage:    "JSON file containing the parsed ingredients and their quantities. Output of Parse command",
			},
			&cli.StringFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "The name of the shopping list the ingredients should be added to",
			},
		},
		Action: func(c *cli.Context) error {
			listName := c.String("list")
			input := c.String("input")
			fmt.Println(input)
			fmt.Println(listName)

			ingredients := make([]model.Ingredient, 0)
			err := ReadJSONFromFile(input, &ingredients)
			if err != nil {
				return err
			}

			err = svc.AddToShoppingList(listName, ingredients)
			if err != nil {
				return err
			}

			if len(ingredients) > 0 {
				fmt.Printf("Successfully added %d ingredients to list %s\n", len(ingredients), listName)
			}

			return nil
		},
	}
}
