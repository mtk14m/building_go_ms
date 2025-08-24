package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"created_on"`
	UpdatedOn   string  `json:"updated_on"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProductById(id)
	if err != nil {
		return err
	}
	//il faut quand même faire la validation du json object reçu
	p.ID = id
	//update the product
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProductById(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}
func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Espresso",
		Description: "Un café corsé préparé avec une petite quantité d'eau chaude sous haute pression.",
		Price:       1.99,
		SKU:         "coffee-esp-001",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          2,
		Name:        "Cappuccino",
		Description: "Un mélange équilibré d'espresso, de lait chaud et de mousse de lait.",
		Price:       2.99,
		SKU:         "coffee-cap-002",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          3,
		Name:        "Latte",
		Description: "Un espresso doux avec une grande quantité de lait chaud et un peu de mousse.",
		Price:       3.49,
		SKU:         "coffee-lat-003",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          4,
		Name:        "Americano",
		Description: "Un espresso allongé avec de l'eau chaude, goût proche du café filtre.",
		Price:       2.49,
		SKU:         "coffee-ame-004",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}
