package dto

type CreateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductInput struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}
