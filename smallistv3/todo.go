package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func createTodoHandler(c *gin.Context) {
	//c.JSON(200, "createTodoHandler") // 这是测试用的
	// 获取请求的参数
	var todo Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "createTodoHandler 参数格式错误",
		})
		return
	}
	// 业务逻辑,往数据库里加一条数据

	// 获取当前用户id
	v, _ := c.Get("uid")
	uid := v.(int64)
	if uid <= 0 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "获取不到uid，请重新登录",
		})
		return
	}
	todo.Uid = uid

	if err := db.Create(&todo).Error; err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "createTodoHandler 服务端异常",
		})
		return
	}
	// 返回响应
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		//"data":todo,
	})
}

func getTodoHandler(c *gin.Context) {
	//c.JSON(200, "getTodoHandler")
	//var todo Todo
	// 请求参数
	// 执行业务逻辑
	var todos []Todo
	if err := db.Find(&todos).Error; err != nil {
		fmt.Println("getTodoHandler查询数据失败", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询数据失败",
		})
		return
	}
	// 返回响应
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": todos,
	})
}

// 这里只更新了状态
func updateTodoHandler(c *gin.Context) {
	//c.JSON(200, "updateTodoHandler")
	// 获取请求参数
	var todo Todo
	if err := c.ShouldBind(&todo); err != nil {
		fmt.Println("updateTodoHandler获取请求参数失败", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数格式错误",
		})
		return
	}
	// 执行业务逻辑
	// 先查数据库是否存在，事实上不需要，因为前端页面数据存在，数据库就存在啊
	if err := db.First(&Todo{}, todo.ID).Error; err != nil {
		// 没查到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有这条记录
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "ErrRecordNotFound",
			})
			return
		}
		// 其他错误
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端异常，请稍候再试",
		})
		return
	}
	// 代码执行到这。说明数据库确实存在todo.ID这条记录
	//
	if err := db.Model(&todo).Update("status", todo.Status).Error; err != nil {
		fmt.Println("updateTodoHandler更新至数据库失败", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端异常，请稍候再试",
		})
		return
	}
	//这样写会有问题的。
	//if err := db.Save(&todo).Error; err != nil {
	//	fmt.Println("updateTodoHandler更新至数据库失败", err)
	//	c.JSON(200, gin.H{
	//		"code": 1,
	//		"msg":  "服务端异常，请稍候再试",
	//	})
	//	return
	//}
	// 返回响应
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}

func deleteTodoHandler(c *gin.Context) {
	//c.JSON(200, "deleteTodoHandler")
	// 请求参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("deleteTodoHandler无效参数", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误",
		})
	}
	// 执行业务落脚
	// 先看一下有没有这条记录
	if err := db.First(&Todo{}, id).Error; err != nil {
		// 没查到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有这条记录
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "ErrRecordNotFound",
			})
			return
		}
		// 其他错误
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端异常，请稍候再试",
		})
		return
	}
	// 删除数据
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端异常，请稍候再试",
		})
		return
	}
	// 返回消息
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
	})
}
