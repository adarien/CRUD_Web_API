package main

import (
	"database/sql"
	"errors"
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

func (db *DB) deleteProduct(c *gin.Context) {
	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		log.Fatal("incorrect ID")
	}
	err = db.deleteProductDB(ID)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{"status": "deleted"})
}

type productForUpdate struct {
	ID    string
	title string
	count string
	price string
}

func (db *DB) checkUpdateProduct(tx *sql.Tx, uP productForUpdate) error {
	var updProduct Product
	ID, err := strconv.ParseInt(uP.ID, 10, 32)
	if err != nil {
		return errors.New("incorrect ID")
	}
	updProduct.ID = ID

	if uP.title != "" {
		updProduct.Title = uP.title
		field := "title"
		err := db.updateProductDB(tx, updProduct, field)
		if err != nil {
			return err
		}
	}

	if uP.count != "" {
		count, err := strconv.ParseInt(uP.count, 10, 32)
		if err != nil {
			return errors.New("incorrect Count")
		}
		updProduct.Count = count
		field := "count"
		err = db.updateProductDB(tx, updProduct, field)
		if err != nil {
			return err
		}
	}

	if uP.price != "" {
		price, err := strconv.ParseFloat(uP.price, 32)
		if err != nil {
			return errors.New("incorrect Price")
		}
		updProduct.Price = price
		field := "price"
		err = db.updateProductDB(tx, updProduct, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) updateProduct(c *gin.Context) {
	tx, err := db.db.Begin()
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback() }()

	var p productForUpdate
	p.ID = c.GetHeader("id")
	p.title = c.GetHeader("title")
	p.count = c.GetHeader("count")
	p.price = c.GetHeader("price")
	if p.title == "" && p.count == "" && p.price == "" {
		log.Panic("empty field")
	}
	err = db.checkUpdateProduct(tx, p)
	if err != nil {
		log.Panic(err)
	}

	_, err = tx.Exec("insert into logs (entity, action) values ($1, $2)",
		"product", "updated")
	if err != nil {
		log.Panic(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
	c.JSON(200, gin.H{"status": "updated"})
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
	err := db.insertProductDB(newProduct)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, newProduct)
}

func (db *DB) insertProductDB(p Product) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

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
	router.PUT("/product", conn.updateProduct)
	router.DELETE("/product", conn.deleteProduct)
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
	defer func() { _ = rows.Close() }()

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
	defer func() { _ = rows.Close() }()

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

func (db *DB) deleteProductDB(id int64) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = db.db.Query("delete from products where id=$1", id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into logs (entity, action) values ($1, $2)",
		"product", "deleted")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) updateProductDB(tx *sql.Tx, p Product, field string) error {
	qwUpdate := fmt.Sprintf("update products set %s=$1 where id=$2", field)
	switch field {
	case "title":
		_, err := tx.Exec(qwUpdate, p.Title, p.ID)
		if err != nil {
			return err
		}
	case "count":
		_, err := tx.Exec(qwUpdate, p.Count, p.ID)
		if err != nil {
			return err
		}
	case "price":
		_, err := tx.Exec(qwUpdate, p.Price, p.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
