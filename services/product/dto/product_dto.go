package dto

// CreateProductRequest represents the payload for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest represents the payload for updating a product
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	CategoryID  *uint    `json:"category_id"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

// ProductResponse represents a product in API responses
type ProductResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	CategoryID  uint              `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	ImageURL    string            `json:"image_url"`
	IsActive    bool              `json:"is_active"`
}

// CreateCategoryRequest represents the payload for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
