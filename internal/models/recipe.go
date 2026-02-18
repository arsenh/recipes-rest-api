package models

import "time"

type Recipe struct {
	ID           string    `json:"id" example:"ckx123abc456"`
	Name         string    `json:"name" example:"Chocolate Chip Cookies" binding:"required"`
	Tags         []string  `json:"tags" example:"[\"dessert\",\"vegan\"]" binding:"required"`
	Ingredients  []string  `json:"ingredients" example:"[\"flour\",\"sugar\",\"chocolate\"]" binding:"required"`
	Instructions []string  `json:"instructions" example:"[\"Mix ingredients\",\"Bake at 180°C\"]" binding:"required"`
	PublishedAt  time.Time `json:"publishedAt" example:"2025-02-18T10:30:00Z"`
}
