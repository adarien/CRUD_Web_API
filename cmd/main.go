package main

import (
	"CRUD_Web_API/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
		logrus.Fatal(err)
	}
}
