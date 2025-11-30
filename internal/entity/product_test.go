package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	prod, err := NewProduct("Xbox", 2900)
	assert.Nil(t, err)
	assert.NotNil(t, prod)
	assert.NotEmpty(t, prod.ID)
	assert.Equal(t, "Xbox", prod.Name)
	assert.Equal(t, 2900, prod.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	prod, err := NewProduct("", 2900)
	assert.Nil(t, prod)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	prod, err := NewProduct("Xbox", -10)
	assert.Nil(t, prod)
	assert.Equal(t, ErrInvalidPrice, err)
}
