package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/adrianomota/fullcycle/my-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getConnection(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	return db
}

func TestCreateProduct(t *testing.T) {
	db := getConnection(t)

	product, err := entity.NewProduct("Product 1", 100.0)
	assert.Nil(t, err)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.Nil(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindALlProducts(t *testing.T) {
	db := getConnection(t)

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.Nil(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)
}

func TestFindProductById(t *testing.T) {
	db := getConnection(t)

	product, err := entity.NewProduct("Product 1", 10.00)
	assert.Nil(t, err)
	db.Create(product)
	productDB := NewProduct(db)
	productExists, err := productDB.FindById(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "Product 1", productExists.Name)
}

func TestUpdateProduct(t *testing.T) {
	db := getConnection(t)

	product, err := entity.NewProduct("Product 1", 10.00)
	assert.Nil(t, err)
	db.Create(product)

	productDB := NewProduct(db)
	product.Name = "Product Updated to 2"
	err = productDB.Update(product)
	assert.Nil(t, err)
	productExists, err := productDB.FindById(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "Product Updated to 2", productExists.Name)
}

func TestRemoveProduct(t *testing.T) {
	db := getConnection(t)

	product, err := entity.NewProduct("Product to remove 1", 10.00)
	assert.Nil(t, err)
	db.Create(product)
	productDB := NewProduct(db)
	productDB.Delete(product.ID.String())

	productExists, err := productDB.FindById(product.ID.String())
	assert.Error(t, err)
	assert.Empty(t, productExists.Name)
}
