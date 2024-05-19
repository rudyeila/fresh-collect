package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/rudyeila/fresh-collect/hellofresh/client"
	"github.com/rudyeila/fresh-collect/service/model"
	"github.com/rudyeila/go-bring-api/bring"
)

type Service struct {
	Bring *bring.Bring
	HF    *client.HelloFresh
	Log   *slog.Logger
}

func (s *Service) AddToShoppingList(listName string, ingredients []model.Ingredient) error {
	lists, err := s.Bring.GetLists()
	if err != nil {
		return err
	}

	if lists == nil {
		return errors.New("no lists found for Bring user")
	}

	listId := ""
	for _, l := range lists.Lists {
		if l.Name == listName {
			listId = l.ListUuid
			break
		}
	}

	if listId == "" {
		return fmt.Errorf("no list with name %s was found for user", listName)
	}

	for _, ingr := range ingredients {
		sub := ""
		if ingr.Amount != nil {
			sub = fmt.Sprintf("%s %s", strconv.FormatFloat(*ingr.Amount, 'f', -1, 64), ingr.Unit)
		}

		err = s.Bring.AddItem(listId, ingr.Name, sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetMergedIngredients(shippedOnly bool, recipeIDs ...string) ([]model.Ingredient, error) {
	var result []model.Ingredient

	recipes, err := s.GetMultipleRecipes(recipeIDs...)
	if err != nil {
		return result, err
	}

	return s.mergeIngredients(shippedOnly, recipes), nil
}

func (s *Service) GetMultipleRecipes(recipeIDs ...string) ([]*model.Recipe, error) {
	result := make([]*model.Recipe, len(recipeIDs))
	for i, id := range recipeIDs {
		recipe, err := s.GetRecipe(id)
		if err != nil {
			return result, err
		}
		result[i] = recipe
	}

	return result, nil
}

func (s *Service) GetRecipe(recipeID string) (*model.Recipe, error) {
	recipe, err := s.HF.GetRecipe(recipeID)
	if err != nil {
		s.Log.Error("getting recipe", "id", recipeID, "error", err.Error())
		return nil, err
	}

	return recipe.ToService(), nil
}

func (s *Service) mergeIngredients(shippedOnly bool, recipes []*model.Recipe) []model.Ingredient {
	ingrMap := make(map[string]model.Ingredient, len(recipes))

	for _, rec := range recipes {
		for _, ingr := range rec.Ingredients {
			if shippedOnly {
				if !ingr.Shipped {
					continue
				}
			}

			merged, ok := ingrMap[ingr.Name]
			if !ok {
				ingrMap[ingr.Name] = ingr
				continue
			}

			// merge new amount with existing one
			if merged.Unit != ingr.Unit {
				s.Log.Warn(fmt.Sprintf("Unit disambiguity for ingredient with Name %s. First unit %s, other unit %s", ingr.Name, ingr.Unit, merged.Unit))
				continue
			}
			if merged.Amount != nil && ingr.Amount != nil {
				*merged.Amount += *ingr.Amount
			}
			// rewrite new value to map
			ingrMap[ingr.Name] = merged
		}
	}

	res := make([]model.Ingredient, 0, len(ingrMap))
	for _, val := range ingrMap {
		res = append(res, val)
	}

	return res
}
