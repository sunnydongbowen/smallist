package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var db *gorm.DB

func initDB() (err error) {
	dsn := "root:815qza@tcp(192.168.72.130:3306)/sql_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	fmt.Printf("connect DB failed,err:%v\n", err)
	//	return
	//}
	return
}

type Todo struct {
	gorm.Model
	Title  string `form:"title" json:"title"` // 代办事项
	Status bool   `json:"status"`             // 是否完成的状态
	Uid    int64  `gorm:"uid;not null;default:0"`
}

// Account 用户表
type Account struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Uid      int64  `gorm:"uid;unique"`  // 用户id 唯一标识
	Name     string `gorm:"name;unique"` // 用户名，不能更改
	Password string `gorm:"password"`
	NickName string `gorm:"nick_name"` // 昵称随便改

	Status *bool `gorm:"status"`
}

func main() {
	if err := initDB(); err != nil {
		fmt.Println("connect mysql err", err)
		panic(err)
	}
	// 映射到数据库里面
	//db.AutoMigrate(&Todo{})
	//db.AutoMigrate(&Account{})

	r := gin.Default()

	// 加载静态文件
	r.LoadHTMLFiles("gin/smallist/index.html")
	r.Static("/static", "gin/smallist/static")

	r.POST("/register", regHandler)
	r.POST("/login", loginHandler)

	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// 增删改查
	g := r.Group("/api/v1", authMiddleware)
	{
		g.POST("/todo", createTodoHandler)
		g.PUT("/todo", updateTodoHandler)
		g.GET("/todo", getTodoHandler)
		g.DELETE("/todo/:id", deleteTodoHandler)

	}

	r.Run()

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

func createTodoHandler(c *gin.Context) {
	//c.JSON(200, "createTodoHandler") // 这是测试用的
	// 获取请求的参数
	var todo Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "无效参数",
		})
		return
	}
	// 业务逻辑
	if err := db.Create(&todo).Error; err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端异常",
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
