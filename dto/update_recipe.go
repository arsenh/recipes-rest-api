package dto

// DTO for create new recipe
type UpdateRecipeRequest struct {
	Name         string   `json:"name" binding:"required"`
	Tags         []string `json:"tags" binding:"required,min=1,dive,required"`
	Ingredients  []string `json:"ingredients" binding:"required,min=1,dive,required"`
	Instructions []string `json:"instructions" binding:"required,min=1,dive,required"`
}
