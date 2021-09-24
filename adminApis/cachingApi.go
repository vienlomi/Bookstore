package adminApis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_v3/data"
	"project_v3/redis"
	"time"

	"github.com/gorilla/mux"
)

func SetListHomeHandler(w http.ResponseWriter, r *http.Request) {
	var input []data.ProductLite
	req := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("decode err", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = redis.ClientRPushProduct(data.OutStanding[req["id"]], input)
	if err != nil {
		log.Println("LPUSH redis fail", err)
		http.Error(w, "push new list fail", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func AddListHomeHandler(w http.ResponseWriter, r *http.Request) {
	var input []data.ProductLite
	req := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("decode err", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = redis.ClientAddRPushProduct(data.OutStanding[req["id"]], input, time.Minute*10)
	if err != nil {
		log.Println("LPUSH redis fail", err)
		http.Error(w, "push new list fail", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func GetListHomeHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	var newProducts []data.ProductLite
	err := redis.ClientLRangProduct(data.OutStanding[req["id"]], &newProducts)
	if err != nil {
		fmt.Println("can not get list")
		http.Error(w, "can not read list redis", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(newProducts)
}

func GetDataHomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get home page")

	var response struct {
		NewProducts []data.ProductLite `json:"new_product"`
		BestSales   []data.ProductLite `json:"best_sales"`
		Popular     []data.ProductLite `json:"popular"`
		HighestRate []data.ProductLite `json:"highest_rate"`
	}

	err := redis.ClientLRangProduct(data.OutStanding["1"], &response.NewProducts)
	if err != nil {
		fmt.Println("can not get list")
		http.Error(w, "can not read list 1 redis", http.StatusInternalServerError)
		return
	}

	err = redis.ClientLRangProduct(data.OutStanding["2"], &response.BestSales)
	if err != nil {
		fmt.Println("can not get list")
		http.Error(w, "can not read list 2 redis", http.StatusInternalServerError)
		return
	}

	err = redis.ClientLRangProduct(data.OutStanding["3"], &response.Popular)
	if err != nil {
		fmt.Println("can not get list")
		http.Error(w, "can not read list 3 redis", http.StatusInternalServerError)
		return
	}

	err = redis.ClientLRangProduct(data.OutStanding["4"], &response.HighestRate)
	if err != nil {
		fmt.Println("can not get list")
		http.Error(w, "can not read list 4 redis", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}
