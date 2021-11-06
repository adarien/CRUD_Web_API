package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Product struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Count   int32     `json:"count"`
	Price   float64   `json:"price"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

type errBR struct {
	Error string
}

var Products = []Product{
	{ID: 1, Title: "Apple", Count: 200, Price: 54.0},
	{ID: 2, Title: "Orange", Count: 250, Price: 72.5},
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Products)
}

func postProduct(c *gin.Context) {
	var newProduct Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, errBR{"bad_request"})
		return
	}

	Products = append(Products, newProduct)
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func main() {
	// router := gin.Default()
	// router.GET("/products", getProducts)
	// router.POST("/products", postProduct)
	//
	// err := router.Run("localhost:8080")
	// if err != nil {
	// 	return
	// }

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres sslmode=disable password=qwerty123")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = insertUser(db, Product{
		Title: "Apple",
		Count: 200,
		Price: 54.0,
	})

	fmt.Println(getProductsSQL(db))
	//
	// if err := db.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	// users, err := getUsers(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(users)
	//
	//     err = insertUser(db, User{
	//        Name: "Петя",
	//        Email: "grisha@ninja.go",
	//     })
	//
	//     users, err = getUsers(db)
	//     if err != nil {
	//        log.Fatal(err)
	//     }
	//     fmt.Println(users)
}

func getProductsSQL(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("select * from products")
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

//
// func getUserByID(db *sql.DB, id int) (User, error) {
// 	var u User
// 	err := db.QueryRow("select * from users where id = $1", id).
// 		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
// 	return u, err
// }
//
func insertUser(db *sql.DB, p Product) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into products (title, count, price) values ($1, $2, $3)",
		p.Title, p.Count, p.Price)
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
