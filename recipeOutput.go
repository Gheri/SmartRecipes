package main

type RecipeOutput struct {
	UniqueRecipeCount       int                     `json:"unique_recipe_count"`
	RecipeVsCount           []RecipeVsCount         `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode         `json:"busiest_postcode"`
	CountPerPostcodeAndTime CountPerPostcodeAndTime `json:"count_per_postcode_and_time"`
	MatchByName             []string                `json:"match_by_name"`
}

type CountPerPostcodeAndTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

type BusiestPostcode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}
type RecipeVsCount struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}
