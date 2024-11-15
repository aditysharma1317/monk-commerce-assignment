package dtos

type Coupon struct {
	Id      string        `json:"id"`
	Type    string        `json:"type"`
	Details CouponDetails `json:"details"`
}

type CouponDetails struct {
	Threshold     int                      `json:"threshold"`
	Discount      int                      `json:"discount"`
	ProductId     string                   `json:"product_id"`
	Quantity      int                      `json:"quantity"`
	BuyProducts   []ProductQuantityDetails `json:"buy_products"`
	GetProducts   []ProductQuantityDetails `json:"get_products"`
	RepitionLimit int                      `json:"repitition_limit"`
}

type ProductQuantityDetails struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Request structure for the POST /applicable-coupons endpoint
type ApplicableCouponsRequest struct {
	Cart Cart `json:"cart"`
}

// Structure representing the shopping cart in the request
type Cart struct {
	Items []CartItem `json:"items"`
}

// Structure representing an item in the cart
type CartItem struct {
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Structure for the response of the POST /applicable-coupons endpoint
type ApplicableCouponsResponse struct {
	ApplicableCoupons []ApplicableCoupon `json:"applicable_coupons"`
}

// Structure representing each applicable coupon in the response
type ApplicableCoupon struct {
	CouponID string  `json:"coupon_id"`
	Type     string  `json:"type"`
	Discount float64 `json:"discount"`
}

type UpdatedCart struct {
	Items         []CartItemDiscount `json:"items"`
	TotalPrice    float64            `json:"total_price"`
	TotalDiscount float64            `json:"total_discount"`
	FinalPrice    float64            `json:"final_price"`
}

type CartItemDiscount struct {
	ProductId     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount"`
}
