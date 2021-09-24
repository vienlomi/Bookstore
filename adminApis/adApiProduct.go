package adminApis

import (
	"encoding/json"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"

	"github.com/gorilla/mux"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Product
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("decode err", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = connection.Model.Product.Insert(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not add new book", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Product
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid product", http.StatusBadRequest)
		return
	}
	err = connection.Model.Product.Update(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not update product", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	err := connection.Model.Product.Delete(id["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

//func AdminSetHomePageHandler(w http.ResponseWriter, r *http.Request) {
//
//}
