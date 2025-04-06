package entity

// ProductData represents a single product record from the dataset
type ProductData struct {
	ProductName         string  `json:"product_name"`
	Brand               string  `json:"brand"`
	Date                string  `json:"date"`
	SalesQuantity       int     `json:"sales_quantity"`
	Price               float64 `json:"price"`
	OriginalPrice       float64 `json:"original_price"`
	DiscountPercentage  float64 `json:"discount_percentage"`
	StockLevel          int     `json:"stock_level"`
	Region              string  `json:"region"`
	Category            string  `json:"category"`
	CustomerRating      float64 `json:"customer_rating"`
	ReviewCount         int     `json:"review_count"`
	DeliveryDays        int     `json:"delivery_days"`
	Seller              string  `json:"seller"`
	IsWeekend           bool    `json:"is_weekend"`
	IsHoliday           bool    `json:"is_holiday"`
}
