package userApis

import (
	"encoding/json"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"
	"project_v3/mail"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/dgrijalva/jwt-go/v4/request"
	"github.com/gorilla/mux"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input data.Order
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = connection.Model.Order.Insert(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not add new order", http.StatusInternalServerError)
		return
	}
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "token expires", http.StatusUnauthorized)
		return
	}
	if token.Valid {
		tokenString := token.Raw
		claims := jwt.MapClaims{}
		token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		})
		// ... error handling
		if err != nil {
			log.Println(err)
			return
		}
		// do something with decoded claims
		var name string
		var email string
		for key, val := range claims {
			if key == "email" {
				email = val.(string)
			}
			if key == "name" {
				name = val.(string)
			}
		}
		mail.Pre.PreSend(name, email, input)
		_ = json.NewEncoder(w).Encode(http.StatusOK)
	} else {
		http.Error(w, "something wrong when sending mail", http.StatusInternalServerError)
	}
}

func GetOrdersOfUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload data.Payload
	_ = json.NewDecoder(r.Body).Decode(&payload)
	if payload.Stt == 0 || payload.Limit == 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	req := mux.Vars(r)
	listOrder, err := connection.Model.Order.GetOrdersByUserId(req["id"], payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "can not get all order", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(listOrder)
}

func GetDetailOrderHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	order, err := connection.Model.Order.GetDetailOrderBuyId(req["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "can not get detail order", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(order)
}
