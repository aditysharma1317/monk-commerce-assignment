package daos

import (
	"monk-commerce-assignment/models"
	"monk-commerce-assignment/utils/context"
)

type Coupon struct {
}

func NewCoupon() ICoupon {
	return &Coupon{}
}

type ICoupon interface {
	PersistCoupon(ctx *context.Context, req *models.Coupon) error
	PersistCartWiseCoupon(ctx *context.Context, req *models.CartWiseCoupon) error
	PersistProductWiseCoupon(ctx *context.Context, req *models.ProductWiseCoupon) error
	PersistBxGyCoupon(ctx *context.Context, req *models.BxGyCoupon) error
	PersistBxGyBuyCoupon(ctx *context.Context, req *models.BxGyBuyProduct) error
	PersistBxGyGetCoupon(ctx *context.Context, req *models.BxGyGetProduct) error
	GetAllCoupons(ctx *context.Context) ([]*models.Coupon, error)
	GetCartWiseCoupon(ctx *context.Context, couponId string) (*models.CartWiseCoupon, error)
	GetProductWiseCoupon(ctx *context.Context, couponId string) (*models.ProductWiseCoupon, error)
	GetBxGyCoupon(ctx *context.Context, couponId string) (*models.BxGyCoupon, error)
	GetBxGyBuyProducts(ctx *context.Context, bxgyCouponId string) ([]*models.BxGyBuyProduct, error)
	GetBxGyGetProducts(ctx *context.Context, bxgyCouponId string) ([]*models.BxGyGetProduct, error)
	GetCouponById(ctx *context.Context, id string) (*models.Coupon, error)
	DeleteCoupon(ctx *context.Context, couponId string) error
	DeleteCartWiseCoupon(ctx *context.Context, couponId string) error
	DeleteProductWiseCoupon(ctx *context.Context, couponId string) error
	DeleteBxGyCoupon(ctx *context.Context, couponId string) error
	DeleteBxGyBuyProducts(ctx *context.Context, couponId string) error
	DeleteBxGyGetProducts(ctx *context.Context, couponId string) error
}

func (c *Coupon) PersistCoupon(ctx *context.Context, req *models.Coupon) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) PersistCartWiseCoupon(ctx *context.Context, req *models.CartWiseCoupon) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) PersistProductWiseCoupon(ctx *context.Context, req *models.ProductWiseCoupon) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) PersistBxGyCoupon(ctx *context.Context, req *models.BxGyCoupon) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) PersistBxGyBuyCoupon(ctx *context.Context, req *models.BxGyBuyProduct) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) PersistBxGyGetCoupon(ctx *context.Context, req *models.BxGyGetProduct) error {
	err := ctx.Transaction.Debug().Create(req).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Coupon) GetAllCoupons(ctx *context.Context) ([]*models.Coupon, error) {
	var coupons []*models.Coupon
	err := ctx.DB.Debug().Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func (c *Coupon) GetCartWiseCoupon(ctx *context.Context, couponId string) (*models.CartWiseCoupon, error) {
	var cartCoupon models.CartWiseCoupon
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).First(&cartCoupon).Error
	if err != nil {
		return nil, err
	}
	return &cartCoupon, nil
}

func (c *Coupon) GetProductWiseCoupon(ctx *context.Context, couponId string) (*models.ProductWiseCoupon, error) {
	var productCoupon models.ProductWiseCoupon
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).First(&productCoupon).Error
	if err != nil {
		return nil, err
	}
	return &productCoupon, nil
}

func (c *Coupon) GetBxGyCoupon(ctx *context.Context, couponId string) (*models.BxGyCoupon, error) {
	var bxgyCoupon models.BxGyCoupon
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).First(&bxgyCoupon).Error
	if err != nil {
		return nil, err
	}
	return &bxgyCoupon, nil
}

func (c *Coupon) GetBxGyBuyProducts(ctx *context.Context, bxgyCouponId string) ([]*models.BxGyBuyProduct, error) {
	var buyProducts []*models.BxGyBuyProduct
	err := ctx.DB.Debug().Where("bxgy_coupon_id = ?", bxgyCouponId).Find(&buyProducts).Error
	if err != nil {
		return nil, err
	}
	return buyProducts, nil
}

func (c *Coupon) GetBxGyGetProducts(ctx *context.Context, bxgyCouponId string) ([]*models.BxGyGetProduct, error) {
	var getProducts []*models.BxGyGetProduct
	err := ctx.DB.Debug().Where("bxgy_coupon_id = ?", bxgyCouponId).Find(&getProducts).Error
	if err != nil {
		return nil, err
	}
	return getProducts, nil
}

func (c *Coupon) GetCouponById(ctx *context.Context, id string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := ctx.DB.Debug().Where("id = ?", id).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (c *Coupon) DeleteCoupon(ctx *context.Context, couponId string) error {
	// Delete the main coupon record
	err := ctx.DB.Debug().Where("id = ?", couponId).Delete(&models.Coupon{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Coupon) DeleteCartWiseCoupon(ctx *context.Context, couponId string) error {
	// Delete the cart-wise coupon entry
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).Delete(&models.CartWiseCoupon{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Coupon) DeleteProductWiseCoupon(ctx *context.Context, couponId string) error {
	// Delete the product-wise coupon entry
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).Delete(&models.ProductWiseCoupon{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Coupon) DeleteBxGyCoupon(ctx *context.Context, couponId string) error {
	// Delete the main BxGy coupon entry
	err := ctx.DB.Debug().Where("coupon_id = ?", couponId).Delete(&models.BxGyCoupon{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Coupon) DeleteBxGyBuyProducts(ctx *context.Context, couponId string) error {
	// Delete BxGy buy products related to the coupon
	err := ctx.DB.Debug().Where("bxgy_coupon_id = ?", couponId).Delete(&models.BxGyBuyProduct{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Coupon) DeleteBxGyGetProducts(ctx *context.Context, couponId string) error {
	// Delete BxGy get products related to the coupon
	err := ctx.DB.Debug().Where("bxgy_coupon_id = ?", couponId).Delete(&models.BxGyGetProduct{}).Error
	if err != nil {
		return err
	}
	return nil
}
