package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {

	p, err := NewProduct("Product 1", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "Product 1", p.Name)
	assert.Equal(t, 10.0, p.Price)
}

func TestProductWhenNameIsrequired(t *testing.T) {
	p, err := NewProduct("", 10.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequire, err)
}

func TestWHenPriceIsrequired(t *testing.T) {
	p, err := NewProduct("Product 1", 0.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestWHenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Product 1", -2.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProductValidate(t *testing.T) {
	p, err := NewProduct("Product 1", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
