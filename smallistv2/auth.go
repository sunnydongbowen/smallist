package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type AuthParam struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册处理函数
func regHandler(c *gin.Context) {
	var param AuthParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "参数错误",
		})
		return
	}
	// 拿到参数去注册用户
	// 去数据库创建一条记录
	var use Account
	err := db.Where("name=?", param.Name).First(&use).Error

	if err == nil {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "用户名已存在",
		})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "服务器异常",
		})
		return
	}
	// 说明没有查到也没有报错
	err = db.Create(&Account{
		Uid:      time.Now().Unix(), // 后面雪花算法？
		Name:     param.Name,
		Password: md5Secret(param.Password),
	}).Error
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "创建用户失败",
		})
		return
	}
	c.JSON(http.StatusOK, Resp{
		Code: 0,
		Msg:  "注册成功",
	})
	return
}

func loginHandler(c *gin.Context) {
	// 1. 获取参数
	var param AuthParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "参数格式错误",
		})
		return
	}
	// 2. 逻辑处理
	var u Account
	err := db.Where("name=? and password=?", param.Name, md5Secret(param.Password)).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, Resp{
				Code: 1,
				Msg:  "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "服务端异常,请稍候再试",
		})
		return
	}
	// 登录成功，接下来发token
	token, err := GenToken(u.Uid, u.Name)
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: 1,
			Msg:  "服务端异常，请售后再试",
		})
		return
	}

	// 3. 返回响应
	c.JSON(http.StatusOK, Resp{
		Code: 0,
		Msg:  "登录成功",
		Data: token,
	})
}

func md5Secret(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd))
	return hex.EncodeToString(h.Sum(nil))
}
