package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var (
	Db *gorm.DB //声明全局
)

// Model
type Todo struct {
	Id int `json:"id"` //和前端传来的参数进行匹配
	Title string `json:"title"`
	Status bool `json:"status"`
}

func initMysql() (err error) {
	info := "root:Pan123456@@tcp(9.135.220.214:3306)/TodoList?charset=utf8&parseTime=True&loc=Local"
	Db, err = gorm.Open("mysql", info)
	if err != nil {
		return
	}
	return Db.DB().Ping()
}

func main() {
	//连接数据库
	err := initMysql()
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	//绑定模型
	Db.AutoMigrate(&Todo{}) //创建表

	r := gin.Default()
	//引入静态文件
	r.Static("/static", "static")
	//告诉go模板页面在哪
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("todo/", func(c *gin.Context) {
			//1. 把前端传到这里的数据取出，也就是绑定数据到后端实体item
			var item Todo
			c.BindJSON(&item)
			//2. 存入数据库
			//3. 返回响应
			if err := Db.Create(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, item) //返回值需要按照前端要求的格式 gin.H{"code": 200, "msg": "success", "data": todo,}
			}
		})

		//查看所有
		v1Group.GET("todo/", func(c *gin.Context) {
			var todoList []Todo
			if err := Db.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusInternalServerError, err.Error)
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})

		//查看单一
		v1Group.GET("todo/:id", func(context *gin.Context) {

		})
		//添加
		v1Group.PUT("todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "DataNotFound"})
				return
			}
			var todo Todo
			if err := Db.Where("id = ?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
				return
			}
			c.BindJSON(&todo)
			if err := Db.Save(todo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//删除
		v1Group.DELETE("todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "DataNotFound"})
				return
			}
			if err := Db.Where("id = ?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"msg": "Success"})
			}
		})
	}

	r.Run()
}