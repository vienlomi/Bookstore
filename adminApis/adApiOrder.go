package adminApis

import (
	"encoding/json"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"
	"project_v3/userApis"
)

func GetDetailOrderHandler(w http.ResponseWriter, r *http.Request) {
	userApis.GetDetailOrderHandler(w, r)
}

func GetOrdersASCHandler(w http.ResponseWriter, r *http.Request) {
	var payload data.Payload
	json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	listOrder, err := connection.Model.Order.GetAllASC(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "can not get all order", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(listOrder)
}

func GetOrdersDESCHandler(w http.ResponseWriter, r *http.Request) {
	var payload data.Payload
	json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	listOrder, err := connection.Model.Order.GetAllDESC(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "can not get all order", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(listOrder)
}
