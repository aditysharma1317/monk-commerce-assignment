package services

import (
	"errors"
	"monk-commerce-assignment/daos"
	"monk-commerce-assignment/dtos"
	"monk-commerce-assignment/models"
	"monk-commerce-assignment/utils/context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CouponService struct {
	db daos.Coupon
}

func NewCouponService() ICouponService {
	return &CouponService{
		db: daos.Coupon{},
	}
}

type ICouponService interface {
	CreateCoupon(ctx *context.Context, req *dtos.Coupon) error
	GetCoupons(ctx *context.Context) ([]*dtos.Coupon, error)
	GetCouponById(ctx *context.Context, id string) (*dtos.Coupon, error)
	GetApplicableCoupons(ctx *context.Context, cartItems []dtos.CartItem) ([]dtos.ApplicableCoupon, error)
	ApplyCoupon(ctx *context.Context, couponId string, cartItems []dtos.CartItem) (*dtos.UpdatedCart, error)
	DeleteCoupon(ctx *context.Context, couponId string) error
}

func (c *CouponService) CreateCoupon(ctx *context.Context, req *dtos.Coupon) error {
	// Start a transaction
	tx := ctx.DB.Begin()
	ctx.Transaction = tx
	if tx.Error != nil {
		ctx.Log.Error("failed to start transaction", zap.Error(tx.Error))
		return tx.Error
	}

	// Generate a new coupon ID and create the base coupon entry
	couponId := uuid.New().String()
	coupon := models.Coupon{
		Id:        couponId,
		Type:      req.Type,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Persist the coupon entry and handle errors with transaction rollback
	err := c.db.PersistCoupon(ctx, &coupon)
	if err != nil {
		ctx.Log.Error("failed to persist coupon", zap.Error(err))
		tx.Rollback()
		return err
	}

	// Process based on coupon type
	switch req.Type {
	case "cart-wise":
		cartWiseCoupon := models.CartWiseCoupon{
			CouponID:  couponId,
			Threshold: float64(req.Details.Threshold),
			Discount:  float64(req.Details.Discount),
		}
		err = c.db.PersistCartWiseCoupon(ctx, &cartWiseCoupon)
		if err != nil {
			ctx.Log.Error("failed to persist cart-wise coupon", zap.Error(err))
			tx.Rollback()
			return err
		}

	case "product-wise":
		productWiseCoupon := models.ProductWiseCoupon{
			CouponID:  couponId,
			ProductID: req.Details.ProductId,
			Discount:  float64(req.Details.Discount),
		}
		err = c.db.PersistProductWiseCoupon(ctx, &productWiseCoupon)
		if err != nil {
			ctx.Log.Error("failed to persist product-wise coupon", zap.Error(err))
			tx.Rollback()
			return err
		}

	case "bxgy":
		bxgyCoupon := models.BxGyCoupon{
			CouponID:        couponId,
			RepetitionLimit: req.Details.RepitionLimit,
		}
		err = c.db.PersistBxGyCoupon(ctx, &bxgyCoupon)
		if err != nil {
			ctx.Log.Error("failed to persist BxGy coupon", zap.Error(err))
			tx.Rollback()
			return err
		}

		// Persist buy products for BxGy coupon
		for _, buyProduct := range req.Details.BuyProducts {
			buyProductModel := models.BxGyBuyProduct{
				BxGyCouponID: couponId,
				ProductID:    buyProduct.ProductId,
				Quantity:     buyProduct.Quantity,
			}
			err = c.db.PersistBxGyBuyCoupon(ctx, &buyProductModel)
			if err != nil {
				ctx.Log.Error("failed to persist BxGy buy product", zap.Error(err))
				tx.Rollback()
				return err
			}
		}

		// Persist get products for BxGy coupon
		for _, getProduct := range req.Details.GetProducts {
			getProductModel := models.BxGyGetProduct{
				BxGyCouponID: couponId,
				ProductID:    getProduct.ProductId,
				Quantity:     getProduct.Quantity,
			}
			err = c.db.PersistBxGyGetCoupon(ctx, &getProductModel)
			if err != nil {
				ctx.Log.Error("failed to persist BxGy get product", zap.Error(err))
				tx.Rollback()
				return err
			}
		}

	default:
		tx.Rollback()
		err := errors.New("unsupported coupon type")
		ctx.Log.Error("invalid coupon type", zap.Error(err))
		return err
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		ctx.Log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (c *CouponService) GetCoupons(ctx *context.Context) ([]*dtos.Coupon, error) {
	// Retrieve all coupons from the database
	coupons, err := c.db.GetAllCoupons(ctx)
	if err != nil {
		return nil, err
	}

	// Prepare a slice to hold the transformed coupon data
	var result []*dtos.Coupon

	// Iterate through each coupon and transform it into the DTO format
	for _, coupon := range coupons {
		couponDto := &dtos.Coupon{
			Id:   coupon.Id,
			Type: coupon.Type,
		}

		// Populate coupon-specific details based on type
		switch coupon.Type {
		case "cart-wise":
			cartCoupon, err := c.db.GetCartWiseCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			couponDto.Details = dtos.CouponDetails{
				Threshold: int(cartCoupon.Threshold),
				Discount:  int(cartCoupon.Discount),
			}

		case "product-wise":
			productCoupon, err := c.db.GetProductWiseCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			couponDto.Details = dtos.CouponDetails{
				ProductId: productCoupon.ProductID,
				Discount:  int(productCoupon.Discount),
			}

		case "bxgy":
			bxgyCoupon, err := c.db.GetBxGyCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			buyProducts, err := c.db.GetBxGyBuyProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}
			getProducts, err := c.db.GetBxGyGetProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}

			// Convert buy and get products into the DTO format
			var buyProductsDto, getProductsDto []dtos.ProductQuantityDetails
			for _, buyProduct := range buyProducts {
				buyProductsDto = append(buyProductsDto, dtos.ProductQuantityDetails{
					ProductId: buyProduct.ProductID,
					Quantity:  buyProduct.Quantity,
				})
			}
			for _, getProduct := range getProducts {
				getProductsDto = append(getProductsDto, dtos.ProductQuantityDetails{
					ProductId: getProduct.ProductID,
					Quantity:  getProduct.Quantity,
				})
			}

			couponDto.Details = dtos.CouponDetails{
				RepitionLimit: bxgyCoupon.RepetitionLimit,
				BuyProducts:   buyProductsDto,
				GetProducts:   getProductsDto,
			}

		default:
			return nil, errors.New("unsupported coupon type")
		}

		// Append the transformed coupon to the result list
		result = append(result, couponDto)
	}

	return result, nil
}

func (c *CouponService) GetCouponById(ctx *context.Context, id string) (*dtos.Coupon, error) {
	// Retrieve the basic coupon information by ID
	coupon, err := c.db.GetCouponById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Initialize the DTO for the response
	couponDto := &dtos.Coupon{
		Id:   coupon.Id,
		Type: coupon.Type,
	}

	// Fetch additional details based on coupon type
	switch coupon.Type {
	case "cart-wise":
		cartCoupon, err := c.db.GetCartWiseCoupon(ctx, coupon.Id)
		if err != nil {
			return nil, err
		}
		couponDto.Details = dtos.CouponDetails{
			Threshold: int(cartCoupon.Threshold),
			Discount:  int(cartCoupon.Discount),
		}

	case "product-wise":
		productCoupon, err := c.db.GetProductWiseCoupon(ctx, coupon.Id)
		if err != nil {
			return nil, err
		}
		couponDto.Details = dtos.CouponDetails{
			ProductId: productCoupon.ProductID,
			Discount:  int(productCoupon.Discount),
		}

	case "bxgy":
		bxgyCoupon, err := c.db.GetBxGyCoupon(ctx, coupon.Id)
		if err != nil {
			return nil, err
		}

		buyProducts, err := c.db.GetBxGyBuyProducts(ctx, bxgyCoupon.CouponID)
		if err != nil {
			return nil, err
		}
		getProducts, err := c.db.GetBxGyGetProducts(ctx, bxgyCoupon.CouponID)
		if err != nil {
			return nil, err
		}

		// Convert buy and get products into the DTO format
		var buyProductsDto, getProductsDto []dtos.ProductQuantityDetails
		for _, buyProduct := range buyProducts {
			buyProductsDto = append(buyProductsDto, dtos.ProductQuantityDetails{
				ProductId: buyProduct.ProductID,
				Quantity:  buyProduct.Quantity,
			})
		}
		for _, getProduct := range getProducts {
			getProductsDto = append(getProductsDto, dtos.ProductQuantityDetails{
				ProductId: getProduct.ProductID,
				Quantity:  getProduct.Quantity,
			})
		}

		couponDto.Details = dtos.CouponDetails{
			RepitionLimit: bxgyCoupon.RepetitionLimit,
			BuyProducts:   buyProductsDto,
			GetProducts:   getProductsDto,
		}

	default:
		return nil, errors.New("unsupported coupon type")
	}

	return couponDto, nil
}

func (c *CouponService) GetApplicableCoupons(ctx *context.Context, cartItems []dtos.CartItem) ([]dtos.ApplicableCoupon, error) {
	// Fetch all available coupons from the database
	coupons, err := c.db.GetAllCoupons(ctx)
	if err != nil {
		return nil, err
	}

	var applicableCoupons []dtos.ApplicableCoupon

	// Calculate the cart total, which may be used by cart-wise coupons
	var cartTotal float64
	for _, item := range cartItems {
		cartTotal += float64(item.Quantity) * item.Price
	}

	// Check each coupon for applicability
	for _, coupon := range coupons {
		var discount float64
		isApplicable := false

		switch coupon.Type {
		case "cart-wise":
			// Check if the cart meets the cart-wise coupon threshold
			cartCoupon, err := c.db.GetCartWiseCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			if cartTotal >= cartCoupon.Threshold {
				discount = (cartCoupon.Discount / 100) * cartTotal
				isApplicable = true
			}

		case "product-wise":
			// Apply discount to specific products in the cart if they match the product-wise coupon
			productCoupon, err := c.db.GetProductWiseCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			for _, item := range cartItems {
				if item.ProductId == productCoupon.ProductID {
					discount += (productCoupon.Discount / 100) * float64(item.Quantity) * item.Price
					isApplicable = true
				}
			}

		case "bxgy":
			// Check BxGy conditions and calculate the discount if applicable
			bxgyCoupon, err := c.db.GetBxGyCoupon(ctx, coupon.Id)
			if err != nil {
				return nil, err
			}
			buyProducts, err := c.db.GetBxGyBuyProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}
			getProducts, err := c.db.GetBxGyGetProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}

			buyCount := 0
			for _, buyProduct := range buyProducts {
				for _, item := range cartItems {
					if item.ProductId == buyProduct.ProductID {
						buyCount += item.Quantity / buyProduct.Quantity
					}
				}
			}

			if buyCount > 0 {
				if buyCount > bxgyCoupon.RepetitionLimit {
					buyCount = bxgyCoupon.RepetitionLimit // Cap buyCount at the repetition limit
				}
				isApplicable = true
				for _, getProduct := range getProducts {
					for _, item := range cartItems {
						if item.ProductId == getProduct.ProductID {
							discount += float64(getProduct.Quantity) * item.Price
						}
					}
				}
				discount *= float64(buyCount)
			}
		}

		// If the coupon is applicable, add it to the result list
		if isApplicable {
			applicableCoupons = append(applicableCoupons, dtos.ApplicableCoupon{
				CouponID: coupon.Id,
				Type:     coupon.Type,
				Discount: discount,
			})
		}
	}

	return applicableCoupons, nil
}

func (c *CouponService) ApplyCoupon(ctx *context.Context, couponId string, cartItems []dtos.CartItem) (*dtos.UpdatedCart, error) {
	// Retrieve the specified coupon by ID
	coupon, err := c.db.GetCouponById(ctx, couponId)
	if err != nil {
		return nil, err
	}

	// Initialize variables for calculating the final prices and discounts
	var totalDiscount float64
	var totalPrice float64
	updatedItems := make([]dtos.CartItemDiscount, len(cartItems))

	// Calculate total cart price
	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.Price
	}

	// Iterate over each item in the cart and calculate discounts based on coupon type
	for i, item := range cartItems {
		updatedItem := dtos.CartItemDiscount{
			ProductId:     item.ProductId,
			Quantity:      item.Quantity,
			Price:         item.Price,
			TotalDiscount: 0,
		}

		switch coupon.Type {
		case "cart-wise":
			// Apply a discount to the entire cart if the threshold is met
			cartCoupon, err := c.db.GetCartWiseCoupon(ctx, couponId)
			if err != nil {
				return nil, err
			}
			if totalPrice >= cartCoupon.Threshold {
				discount := (cartCoupon.Discount / 100) * totalPrice
				totalDiscount += discount
			}

		case "product-wise":
			// Apply a discount to specific products in the cart
			productCoupon, err := c.db.GetProductWiseCoupon(ctx, couponId)
			if err != nil {
				return nil, err
			}
			if item.ProductId == productCoupon.ProductID {
				discount := (productCoupon.Discount / 100) * float64(item.Quantity) * item.Price
				updatedItem.TotalDiscount = discount
				totalDiscount += discount
			}

		case "bxgy":
			// Apply "Buy X, Get Y" coupon logic
			bxgyCoupon, err := c.db.GetBxGyCoupon(ctx, couponId)
			if err != nil {
				return nil, err
			}
			buyProducts, err := c.db.GetBxGyBuyProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}
			getProducts, err := c.db.GetBxGyGetProducts(ctx, bxgyCoupon.CouponID)
			if err != nil {
				return nil, err
			}

			buyCount := 0
			for _, buyProduct := range buyProducts {
				if item.ProductId == buyProduct.ProductID {
					buyCount += item.Quantity / buyProduct.Quantity
				}
			}

			if buyCount > 0 && buyCount <= bxgyCoupon.RepetitionLimit {
				for _, getProduct := range getProducts {
					for _, cartItem := range cartItems {
						if cartItem.ProductId == getProduct.ProductID {
							discount := float64(getProduct.Quantity*buyCount) * cartItem.Price
							updatedItem.TotalDiscount += discount
							totalDiscount += discount
						}
					}
				}
			}
		}

		// Add the updated item to the list of updated items
		updatedItems[i] = updatedItem
	}

	// Calculate final price after applying total discount
	finalPrice := totalPrice - totalDiscount

	// Create the updated cart response
	updatedCart := &dtos.UpdatedCart{
		Items:         updatedItems,
		TotalPrice:    totalPrice,
		TotalDiscount: totalDiscount,
		FinalPrice:    finalPrice,
	}

	return updatedCart, nil
}

func (c *CouponService) DeleteCoupon(ctx *context.Context, couponId string) error {
	// Check if the coupon exists
	coupon, err := c.db.GetCouponById(ctx, couponId)
	if err != nil {
		return errors.New("coupon not found")
	}

	tx := ctx.DB.Begin()
	ctx.Transaction = tx
	if tx.Error != nil {
		ctx.Log.Error("failed to start transaction", zap.Error(tx.Error))
		return tx.Error
	}

	// Delete the coupon based on its type (if additional tables are used for specific types)
	switch coupon.Type {
	case "cart-wise":
		err = c.db.DeleteCartWiseCoupon(ctx, couponId)
	case "product-wise":
		err = c.db.DeleteProductWiseCoupon(ctx, couponId)
	case "bxgy":
		err = c.db.DeleteBxGyCoupon(ctx, couponId)
		if err != nil {
			ctx.Log.Error("error deleting bxgy coupon", zap.Error(err))
			tx.Rollback()
			return err
		}
		err = c.db.DeleteBxGyBuyProducts(ctx, couponId)
		if err != nil {
			ctx.Log.Error("error deleting bxgy buy coupon", zap.Error(err))
			tx.Rollback()
			return err
		}
		err = c.db.DeleteBxGyGetProducts(ctx, couponId)
		if err != nil {
			ctx.Log.Error("error deleting bxgy buy coupon", zap.Error(err))
			tx.Rollback()
			return err
		}
	}
	if err != nil {
		return err
	}

	// Delete the main coupon entry
	err = c.db.DeleteCoupon(ctx, couponId)
	if err != nil {
		ctx.Log.Error("error deleting main coupon", zap.Error(err))
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		ctx.Log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}
