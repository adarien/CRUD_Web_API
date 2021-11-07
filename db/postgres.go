package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

type DB struct {
	DB *sql.DB
}

type Product struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Count   int64     `json:"count"`
	Price   float64   `json:"price"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

func Connect() *DB {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	driverName := viper.GetString("DRIVER")
	host := viper.GetString("HOST")
	port := viper.GetString("PORT")
	user := viper.GetString("USER")
	dbname := viper.GetString("DBNAME")
	sslMode := viper.GetString("SSLMODE")
	password := viper.GetString("PASSWORD")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		host, port, user, dbname, sslMode, password)
	fmt.Println(dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{DB: db}
}

func (db *DB) InsertProductDB(p Product) error {
	tx, err := db.DB.Begin()
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

func (db *DB) GetProductsDB() ([]Product, error) {
	rows, err := db.DB.Query("select * from products")
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

func (db *DB) GetProductDB(id int64) (Product, error) {
	product := Product{}
	rows, err := db.DB.Query("select * from products where id=$1", id)
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

func (db *DB) DeleteProductDB(id int64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = db.DB.Query("delete from products where id=$1", id)
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

func (db *DB) UpdateProductDB(tx *sql.Tx, p Product, field string) error {
	qwUpdate := fmt.Sprintf("update products set %s=$1, updated=now() where id=$2", field)
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
