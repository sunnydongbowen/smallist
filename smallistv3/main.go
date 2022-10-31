package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
