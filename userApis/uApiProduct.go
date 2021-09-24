package userApis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"
	"project_v3/redis"
	"time"

	"github.com/gorilla/mux"
)

func GetDetailProductHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	product, err := connection.Model.Product.GetDetailById(req["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(&product)
}

func GetAllProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload data.Payload
	_ = json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	products, err := connection.Model.Product.GetAll(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not get products", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(products)
}

// show more in homepage
func GetListProductHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	orderBy := data.OrderBy[req["order-by"]]

	var products []data.ProductLite

	var payload data.Payload
	_ = json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 || orderBy == "" {
		http.Error(w, "invalid payload or order by", http.StatusBadRequest)
		return
	}

	err := redis.ClientLRangProduct(orderBy, &products)
	if err == nil {
		fmt.Println("get list in redis success")
		_ = json.NewEncoder(w).Encode(products)
		return
	}
	fmt.Println("redis ko co")

	products, err = connection.Model.Product.GetListProduct(payload, orderBy)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not get products", http.StatusInternalServerError)
		return
	}

	err = redis.ClientRPushProductExpire(orderBy, products, time.Second*100)
	if err != nil {
		fmt.Println("can not set list again in redis")
	}

	fmt.Println("set list again in redis: ", orderBy)

	_ = json.NewEncoder(w).Encode(products)
}
