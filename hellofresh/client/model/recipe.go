package model

import (
	"time"

	"github.com/rudyeila/hello-fresh-go-client/hellofresh/service/model"
)

type Allergen struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Slug             string  `json:"slug"`
	IconPath         *string `json:"iconPath"`
	TriggersTracesOf bool    `json:"triggersTracesOf"`
	TracesOf         bool    `json:"tracesOf"`
	IconLink         string  `json:"iconLink,omitempty"`
}

type Cuisine struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IconLink string `json:"iconLink"`
}

type IngredientFamily struct {
	Id       string      `json:"id"`
	Uuid     string      `json:"uuid"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Type     string      `json:"type"`
	Priority int         `json:"priority"`
	IconLink interface{} `json:"iconLink"`
	IconPath interface{} `json:"iconPath"`
}

type Ingredient struct {
	Id        string           `json:"id"`
	Uuid      string           `json:"uuid"`
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Slug      string           `json:"slug"`
	Country   string           `json:"country"`
	ImageLink string           `json:"imageLink"`
	ImagePath string           `json:"imagePath"`
	Shipped   bool             `json:"shipped"`
	Allergens []string         `json:"allergens"`
	Family    IngredientFamily `json:"family"`
}

type Nutrition struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Tag struct {
	Id           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	ColorHandle  string   `json:"colorHandle"`
	Preferences  []string `json:"preferences"`
	DisplayLabel bool     `json:"displayLabel"`
}

type Yield struct {
	Yields      int `json:"yields"`
	Ingredients []struct {
		Id     string   `json:"id"`
		Amount *float64 `json:"amount"`
		Unit   string   `json:"unit"`
	} `json:"ingredients"`
}

type Step struct {
	Index                int           `json:"index"`
	Instructions         string        `json:"instructions"`
	InstructionsHTML     string        `json:"instructionsHTML"`
	InstructionsMarkdown string        `json:"instructionsMarkdown"`
	Ingredients          []interface{} `json:"ingredients"`
	Utensils             []string      `json:"utensils"`
	Timers               []interface{} `json:"timers"`
	Images               []struct {
		Link    string `json:"link"`
		Path    string `json:"path"`
		Caption string `json:"caption"`
	} `json:"images"`
	Videos []interface{} `json:"videos"`
}

type GetRecipeResponse struct {
	Active              bool         `json:"active"`
	Allergens           []Allergen   `json:"allergens"`
	AverageRating       float64      `json:"averageRating"`
	Canonical           string       `json:"canonical"`
	CanonicalLink       string       `json:"canonicalLink"`
	CardLink            string       `json:"cardLink"`
	Category            interface{}  `json:"category"`
	ClonedFrom          string       `json:"clonedFrom"`
	Comment             interface{}  `json:"comment"`
	Country             string       `json:"country"`
	CreatedAt           time.Time    `json:"createdAt"`
	Cuisines            []Cuisine    `json:"cuisines"`
	Description         string       `json:"description"`
	DescriptionHTML     string       `json:"descriptionHTML"`
	DescriptionMarkdown string       `json:"descriptionMarkdown"`
	Difficulty          int          `json:"difficulty"`
	FavoritesCount      int          `json:"favoritesCount"`
	Headline            string       `json:"headline"`
	Id                  string       `json:"id"`
	ImageLink           string       `json:"imageLink"`
	ImagePath           string       `json:"imagePath"`
	Ingredients         []Ingredient `json:"ingredients"`
	IsAddon             bool         `json:"isAddon"`
	IsComplete          *bool        `json:"isComplete"`
	Link                string       `json:"link"`
	Name                string       `json:"name"`
	Nutrition           []Nutrition  `json:"nutrition"`
	PrepTime            string       `json:"prepTime"`
	RatingsCount        int          `json:"ratingsCount"`
	ServingSize         int          `json:"servingSize"`
	Slug                string       `json:"slug"`
	Steps               []Step       `json:"steps"`
	Tags                []Tag        `json:"tags"`
	TotalTime           string       `json:"totalTime"`
	UniqueRecipeCode    string       `json:"uniqueRecipeCode"`
	UpdatedAt           time.Time    `json:"updatedAt"`
	Uuid                string       `json:"uuid"`
	VideoLink           *string      `json:"videoLink"`
	WebsiteUrl          string       `json:"websiteUrl"`
	Yields              []Yield      `json:"yields"`
}

func (res *GetRecipeResponse) ToService() *model.Recipe {
	var yieldForTwo Yield
	for _, y := range res.Yields {
		if y.Yields == 2 {
			yieldForTwo = y
			break
		}
	}

	ingredientsWithAmount := make([]model.Ingredient, len(res.Ingredients))
	for i, ingr := range res.Ingredients {
		for _, yIngr := range yieldForTwo.Ingredients {
			if ingr.Id == yIngr.Id {
				ingredientsWithAmount[i] = model.Ingredient{
					Id:      ingr.Uuid,
					Name:    ingr.Name,
					Amount:  yIngr.Amount,
					Unit:    yIngr.Unit,
					Shipped: ingr.Shipped,
				}
			}
		}
	}

	return &model.Recipe{ID: res.Uuid, Country: res.Country, Ingredients: ingredientsWithAmount}
}
