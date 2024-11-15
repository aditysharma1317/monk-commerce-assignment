package models

import (
	"time"
)

type Coupon struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	Type      string    `json:"type"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartWiseCoupon struct {
	CouponID  string  `gorm:"primaryKey"`
	Threshold float64 `json:"threshold"`
	Discount  float64 `json:"discount"`
}

type ProductWiseCoupon struct {
	CouponID  string  `gorm:"primaryKey"`
	ProductID string  `json:"product_id"`
	Discount  float64 `json:"discount"`
}

type BxGyCoupon struct {
	CouponID        string `gorm:"primaryKey"`
	RepetitionLimit int    `json:"repetition_limit"`
}

type BxGyBuyProduct struct {
	BxGyCouponID string `gorm:"primaryKey"`
	ProductID    string `json:"product_id"`
	Quantity     int    `json:"quantity"`
}

type BxGyGetProduct struct {
	BxGyCouponID string `gorm:"primaryKey"`
	ProductID    string `json:"product_id"`
	Quantity     int    `json:"quantity"`
}
