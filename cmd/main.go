package main

import (
	l "CRUD_Web_API/logs"
	"CRUD_Web_API/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	conn := service.New()
	router := gin.Default()
	router.GET("/products", conn.GetProducts)
	router.GET("/product", conn.GetProduct)
	router.POST("/products", conn.PostProduct)
	router.PUT("/product", conn.UpdateProduct)
	router.DELETE("/product", conn.DeleteProduct)

	err := router.Run("localhost:8080")
	if err != nil {
		l.ERROR.Fatal(err)
	}
}
