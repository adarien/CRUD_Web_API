package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

type Product struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Count   int64     `json:"count"`
	Price   float64   `json:"price"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

type DB struct {
	db *sql.DB
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

func Connect() *DB {
	driverName, _ := os.LookupEnv("DRIVER")
	host, _ := os.LookupEnv("HOST")
	port, _ := os.LookupEnv("PORT")
	user, _ := os.LookupEnv("USER")
	dbname, _ := os.LookupEnv("DBNAME")
	sslMode, _ := os.LookupEnv("SSLMODE")
	password, _ := os.LookupEnv("PASSWORD")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		host, port, user, dbname, sslMode, password)
	fmt.Println(dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{db: db}
}

func (db *DB) getProducts(c *gin.Context) {
	p, err := db.getProductsDB()
	if err != nil {
		return
	}
	c.JSON(200, p)
}

func (db *DB) getProduct(c *gin.Context) {
	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		log.Fatal("incorrect ID")
	}
	p, err := db.getProductDB(ID)
	if err != nil {
		return
	}
	c.JSON(200, p)
}

func parseNewProduct(c *gin.Context) Product {
	var newProduct Product

	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		log.Fatal("incorrect ID")
	}

	title := c.GetHeader("title")
	if len(title) == 0 {
		log.Fatal("incorrect Title")
	}

	count, err := strconv.ParseInt(c.GetHeader("count"), 10, 32)
	if err != nil {
		log.Fatal("incorrect Count")
	}

	price, err := strconv.ParseFloat(c.GetHeader("price"), 32)
	if err != nil {
		log.Fatal("incorrect Price")
	}

	newProduct.ID = ID
	newProduct.Title = title
	newProduct.Price = price
	newProduct.Count = count
	newProduct.Created = time.Now()
	newProduct.Updated = time.Now()

	return newProduct
}

func (db *DB) postProduct(c *gin.Context) {
	newProduct := parseNewProduct(c)
	err := db.insertProduct(newProduct)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, newProduct)
}

func (db *DB) insertProduct(p Product) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into products (id, title, count, price) values ($1, $2, $3, $4)",
		p.ID, p.Title, p.Count, p.Price)
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into logs (entity, action) values ($1, $2)",
		"product", "created")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func main() {
	conn := Connect()
	router := gin.Default()
	router.GET("/products", conn.getProducts)
	router.GET("/product", conn.getProduct)
	router.POST("/products", conn.postProduct)
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) getProductsDB() ([]Product, error) {
	rows, err := db.db.Query("select * from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.ID, &p.Title, &p.Count, &p.Price, &p.Created, &p.Updated)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (db *DB) getProductDB(id int64) (Product, error) {
	product := Product{}
	rows, err := db.db.Query("select * from products where id=$1", id)
	if err != nil {
		return product, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Title, &product.Count,
			&product.Price, &product.Created, &product.Updated)
		if err != nil {
			return product, err
		}
	}

	err = rows.Err()
	if err != nil {
		return product, err
	}

	return product, nil
}

//
// func deleteUser(db *sql.DB, id int) error {
// 	_, err := db.Exec("delete from users where id = $1", id)
//
// 	return err
// }
//
// func updateUser(db *sql.DB, id int, newUser User) error {
// 	_, err := db.Exec("update users set name=$1, email=$2 where id=$3",
// 		newUser.Name, newUser.Email, id)
//
// 	return err
// }
