package domain

// CartItem represents an item in the shopping cart
type CartItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID     uint       `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPrice float64    `json:"total_price"`
}

// CalculateTotals calculates total items and price
func (c *Cart) CalculateTotals() {
	c.TotalItems = 0
	c.TotalPrice = 0
	for _, item := range c.Items {
		c.TotalItems += item.Quantity
		c.TotalPrice += item.Price * float64(item.Quantity)
	}
}

// AddItem adds an item to the cart or updates quantity if exists
func (c *Cart) AddItem(item CartItem) {
	for i, existing := range c.Items {
		if existing.ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			c.CalculateTotals()
			return
		}
	}
	c.Items = append(c.Items, item)
	c.CalculateTotals()
}

// UpdateItemQuantity updates the quantity of an item
func (c *Cart) UpdateItemQuantity(productID uint, quantity int) bool {
	for i, item := range c.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			} else {
				c.Items[i].Quantity = quantity
			}
			c.CalculateTotals()
			return true
		}
	}
	return false
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID uint) bool {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.CalculateTotals()
			return true
		}
	}
	return false
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.TotalItems = 0
	c.TotalPrice = 0
}
