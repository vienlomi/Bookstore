package adminApis

import (
	"encoding/json"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"

	"github.com/gorilla/mux"
)

func GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload data.Payload
	_ = json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	users, err := connection.Model.User.GetAllUser(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "can't get all users", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(&users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	err := connection.Model.User.Delete(req["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid id user", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}
