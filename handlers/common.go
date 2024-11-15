package handlers

import (
	"monk-commerce-assignment/config"
	"monk-commerce-assignment/utils/context"
	"monk-commerce-assignment/utils/db"
	"monk-commerce-assignment/utils/log"

	"github.com/google/uuid" // Import the uuid package
)

func logAndGetContext(c *context.Context) {
	c.RefID = c.Request.Header.Get("X-Request-Id")

	if c.RefID == "" {
		c.RefID = uuid.New().String() // Use uuid.New() instead of uuid.New().String()
	}

	cfg := config.Get()

	c.Log = log.New(c.RefID, cfg.AppName, cfg.LogLevel)
	c.DB = db.New()
}
