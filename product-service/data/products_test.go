package data

import (
	"testing"
	"time"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:      "minou",
		Price:     1.00,
		SKU:       "coffee-min-001",
		CreatedOn: time.Now().String(),
		UpdatedOn: time.Now().String(),
	}

	err := p.ValidateProduct()
	if err != nil {
		t.Fatal(err)
	}
}
