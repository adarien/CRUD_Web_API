package service

import (
	"CRUD_Web_API/db"
	l "CRUD_Web_API/logs"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type Service struct {
	db *db.DB
}

func New() *Service {
	dbClient := db.Connect()
	return &Service{db: dbClient}
}

func parseNewProduct(c *gin.Context) db.Product {
	var newProduct db.Product

	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		l.ERROR.Fatal("incorrect ID")
	}

	title := c.GetHeader("title")
	if len(title) == 0 {
		l.ERROR.Fatal("incorrect Title")
	}

	count, err := strconv.ParseInt(c.GetHeader("count"), 10, 32)
	if err != nil {
		l.ERROR.Fatal("incorrect Count")
	}

	price, err := strconv.ParseFloat(c.GetHeader("price"), 32)
	if err != nil {
		l.ERROR.Fatal("incorrect Price")
	}

	newProduct.ID = ID
	newProduct.Title = title
	newProduct.Price = price
	newProduct.Count = count
	newProduct.Created = time.Now()
	newProduct.Updated = time.Now()

	return newProduct
}

func (s *Service) PostProduct(c *gin.Context) {
	newProduct := parseNewProduct(c)
	err := s.db.InsertProductDB(newProduct)
	if err != nil {
		l.ERROR.Fatal(err)
	}
	c.JSON(200, newProduct)
}

func (s *Service) GetProducts(c *gin.Context) {
	p, err := s.db.GetProductsDB()
	if err != nil {
		return
	}
	c.JSON(200, p)
}

func (s *Service) GetProduct(c *gin.Context) {
	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		l.ERROR.Fatal("incorrect ID")
	}
	p, err := s.db.GetProductDB(ID)
	if err != nil {
		l.INFO.Panic(err)
	}
	c.JSON(200, p)
}

func (s *Service) DeleteProduct(c *gin.Context) {
	ID, err := strconv.ParseInt(c.GetHeader("id"), 10, 32)
	if err != nil {
		l.ERROR.Fatal("incorrect ID")
	}
	err = s.db.DeleteProductDB(ID)
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

func (s *Service) checkUpdateProduct(tx *sql.Tx, uP productForUpdate) error {
	var updProduct db.Product
	ID, err := strconv.ParseInt(uP.ID, 10, 32)
	if err != nil {
		return errors.New("incorrect ID")
	}
	updProduct.ID = ID

	if uP.title != "" {
		updProduct.Title = uP.title
		field := "title"
		err := s.db.UpdateProductDB(tx, updProduct, field)
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
		err = s.db.UpdateProductDB(tx, updProduct, field)
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
		err = s.db.UpdateProductDB(tx, updProduct, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) UpdateProduct(c *gin.Context) {
	tx, err := s.db.DB.Begin()
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
		l.ERROR.Panic("empty field")
	}
	err = s.checkUpdateProduct(tx, p)
	if err != nil {
		l.ERROR.Panic(err)
	}

	_, err = tx.Exec("insert into logs (entity, action) values ($1, $2)",
		"product", "updated")
	if err != nil {
		l.ERROR.Panic(err)
	}
	err = tx.Commit()
	if err != nil {
		l.ERROR.Panic(err)
	}
	c.JSON(200, gin.H{"status": "updated"})
}
