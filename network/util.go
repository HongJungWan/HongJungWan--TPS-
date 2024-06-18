// Helper Method 역할

package network

import (
	"github.com/gin-gonic/gin"
	"go_chat/types"
)

func Response(c *gin.Context, service int, res interface{}, data ...string) {
	c.JSON(service, types.NewRes(service, res, data...))
}
