package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/ctxkey"
)

func ResponseID() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(ctxkey.ResponseAPI, true)
		c.Next()
		if !common.IsRedisEnabled() {
			return
		}
		id := c.GetString(ctxkey.ResponseAPIID)
		if id == "" {
			return
		}
		userID := c.GetInt(ctxkey.Id)
		if userID == 0 {
			return
		}
		rdb := common.RDB
		key := fmt.Sprintf("response:api:id:%s:%d", id, userID)
		rdb.SetEX(c, key, c.GetString(ctxkey.RequestModel), 10*24*time.Hour)
	}
}

func SetModelIDForRequestAPI() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Param(ctxkey.ResponseAPIID) == "" {
			c.Next()
			return
		}
		if c.GetString(ctxkey.RequestModel) != "" {
			c.Next()
			return
		}
		if !common.IsRedisEnabled() {
			c.Next()
			return
		}
		userID := c.GetInt(ctxkey.Id)
		if userID == 0 {
			c.Next()
			return
		}
		rdb := common.RDB
		key := fmt.Sprintf("response:api:id:%s:%d", c.Param(ctxkey.ResponseAPIID), userID)
		result, err := rdb.Get(c, key).Result()
		if err != nil {
			c.Next()
			return
		}
		c.Set(ctxkey.RequestModel, result)
	}
}
