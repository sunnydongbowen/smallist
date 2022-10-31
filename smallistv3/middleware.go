package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var ctxKey string

// var(
//
//	KeyUid
//
// )
// 中间件
func authMiddleware(c *gin.Context) {
	//1. 从请求头获取token
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Println("authHeader: " + authHeader)

	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "请求头中Bearer Token为空",
		})
		c.Abort()
		return
	}

	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "请求头中Bearer Auth格式有误",
		})
		c.Abort()
		return
	}

	// token
	mc, err := ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "无效的token",
			"err":  err,
		})
		fmt.Println(err)
		c.Abort()
		return
	}
	c.Set("name", mc.Name)
	c.Set("uid", mc.Uid) // 这个uid是拿到的
}
