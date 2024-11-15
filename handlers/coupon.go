package handlers

import (
	"encoding/json"
	"monk-commerce-assignment/dtos"
	"monk-commerce-assignment/services"
	"monk-commerce-assignment/utils/context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/coupons", createCoupon)
	router.GET("/coupons", getAllCoupons)
	router.GET("/coupons/:id", getCouponById)
	router.POST("/applicable-coupons", getApplicableCoupons)
	router.POST("/apply-coupon/:id", applyCoupon)
	router.DELETE("/coupons/:id", deleteCoupon)

}

func createCoupon(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}
	logAndGetContext(ctx)
	decoder := json.NewDecoder(ctx.Request.Body)

	req := &dtos.Coupon{}
	err := decoder.Decode(req)
	if err != nil {
		ctx.Log.Error("error parsing request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	err = services.NewCouponService().CreateCoupon(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string("success"),
	})
}

func getAllCoupons(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	coupons, err := services.NewCouponService().GetCoupons(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

func getCouponById(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	couponId := c.Param("id")

	coupon, err := services.NewCouponService().GetCouponById(ctx, couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

func getApplicableCoupons(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	var request dtos.ApplicableCouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	applicableCoupons, err := services.NewCouponService().GetApplicableCoupons(ctx, request.Cart.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dtos.ApplicableCouponsResponse{
		ApplicableCoupons: applicableCoupons,
	}

	c.JSON(http.StatusOK, response)
}

func applyCoupon(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	couponId := c.Param("id")

	var request dtos.ApplicableCouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	updatedCart, err := services.NewCouponService().ApplyCoupon(ctx, couponId, request.Cart.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

func deleteCoupon(c *gin.Context) {
	// Initialize the context
	ctx := &context.Context{
		Context: c,
	}

	// Get the coupon ID from the URL parameter
	couponId := c.Param("id")

	// Call the service to delete the coupon
	err := services.NewCouponService().DeleteCoupon(ctx, couponId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon deleted successfully",
	})
}
