package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	router := gin.Default()
	router.GET("/products", getProducts)
	router.POST("/products", postProduct)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}

	// db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres sslmode=disable password=qwerty123")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	//
	// if err := db.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	// err = insertUser(db, User{
	//     Name:     "Petya",
	//     Email:    "petya@gmail.com",
	//     Password: "kjrhg754yv754yck545g45h54h55gg",
	// })

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

// func getUsers(db *sql.DB) ([]User, error) {
// 	rows, err := db.Query("select * from users")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	users := make([]User, 0)
// 	for rows.Next() {
// 		u := User{}
// 		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		users = append(users, u)
// 	}
//
// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return users, nil
// }
//
// func getUserByID(db *sql.DB, id int) (User, error) {
// 	var u User
// 	err := db.QueryRow("select * from users where id = $1", id).
// 		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
// 	return u, err
// }
//
// func insertUser(db *sql.DB, u User) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()
//
// 	_, err = tx.Exec("insert into users (name, email, password) values ($1, $2, $3)",
// 		u.Name, u.Email, u.Password)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = tx.Exec("insert into logs (entity, action) values ($1, $2)",
// 		"user", "created")
// 	if err != nil {
// 		return err
// 	}
//
// 	return tx.Commit()
// }
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
