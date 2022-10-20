package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	db := connect()

	defer db.Close()

	//product := Product{ID: uuid.New(), Name: "New Mission Impossible", Quantity: 18, Price: 87.99}
	//fmt.Println(product)

	//fmt.Println(uuid.New())

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/products", createProduct).Methods("POST")
	r.HandleFunc("/api/v1/products", getProducts).Methods("GET")
	r.HandleFunc("/api/v1/product/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/v1/product/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/api/v1/product/{id}", updateProduct).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8090", r))
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := connect()
	defer db.Close()

	product := &Product{
		ID: uuid.New().String(),
	}

	_ = json.NewDecoder(r.Body).Decode(&product)

	_, err := db.Model(product).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(product)
}

//Get Product

func getProducts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//Database connect
	db := connect()
	defer db.Close()

	//Crating Product
	var products []Product
	if err := db.Model(&products).Select(); err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)

	}

	//Returning Objects
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get DB to Connect
	db := connect()
	defer db.Close()

	//Get Id of Product
	params := mux.Vars(r)
	productId := params["id"]

	product := &Product{ID: productId}
	if err := db.Model(product).WherePK().Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Return the Object with the ID passed

	json.NewEncoder(w).Encode(product)
}

//Delete Product

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Database connect

	db := connect()
	defer db.Close()

	//Get ID

	params := mux.Vars(r)
	productId := params["id"]

	//Creating delete query instance for Database
	product := &Product{}
	result, err := db.Model(product).Where("id = ?", productId).Delete()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Return result
	json.NewEncoder(w).Encode(result)
}

//Update Product

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//DB Connect Instace
	db := connect()
	defer db.Close()

	//Get Product id
	params := mux.Vars(r)
	productId := params["id"]

	//Creating product Update instance

	product := &Product{ID: productId}
	_ = json.NewDecoder(r.Body).Decode(&product)

	_, err := db.Model(product).WherePK().Set("name = ?, quantity = ?, price = ?, store = ?", product.Name, product.Quantity, product.Price, product.Store).Update()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Return Updated Product

	json.NewEncoder(w).Encode(product)
}
