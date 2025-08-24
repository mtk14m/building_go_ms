package handlers

import (
	"github/mtk14minou/product-service/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//getProducts
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	//addProduct
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	//updateProduct
	if r.Method == http.MethodPut {
		//path
		path := r.URL.Path
		//getting the id from url
		idRegex := regexp.MustCompile(`/([0-9]+)`)
		g := idRegex.FindStringSubmatch(path)

		if len(g) != 2 {
			p.l.Println("Invalid URI more than one id", g)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[1]

		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI, unable to convert to number", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, w, r)

	}

	//catch all
	w.WriteHeader(http.StatusNotImplemented)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {

	listOfProducts := data.GetProducts()
	err := listOfProducts.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
	p.l.Println("Get all products")
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	product := &data.Product{}
	//should add a json validation before
	//storing product in productList
	data.AddProduct(product)

	//loging product
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshall json", http.StatusBadRequest)
	}
	p.l.Printf("product: %#v", product)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")

	product := &data.Product{}
	//should add a json validation before

	//decode the json object
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	//force the id
	product.ID = id

	//update the product
	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//loging product
	p.l.Printf("product: %#v", product)
}
