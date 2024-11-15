package context

import (
	"monk-commerce-assignment/utils/db"
	"monk-commerce-assignment/utils/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Context struct {
	Log          log.Logger `json:"log"`
	DB           *db.DBConn `json:"db"`
	RefID        string     `json:"ref_id"`
	Transaction  *gorm.DB
	*gin.Context `json:"context"`
}

func (c *Context) Copy() *Context {
	return &Context{
		Log:         c.Log,
		DB:          c.DB,
		RefID:       c.RefID,
		Transaction: c.Transaction,
		Context:     c.Context.Copy(),
	}
}
