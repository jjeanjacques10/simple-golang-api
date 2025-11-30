package database

import (
	"apis/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Laptop", 1500)
	productDB := NewProductDB(db)
	err = productDB.Create(product)
	assert.Nil(t, err)
	var productFound entity.Product
	err = db.First(&productFound, "name = ?", "Laptop").Error
	assert.Nil(t, err)
	assert.Equal(t, "Laptop", productFound.Name)
	assert.Equal(t, 1500, productFound.Price)
}
